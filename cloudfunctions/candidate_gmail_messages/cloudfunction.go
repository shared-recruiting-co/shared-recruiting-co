package cloudfunctions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/getsentry/sentry-go"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const (
	provider = "google"
)

func init() {
	functions.CloudEvent("Handler", handler)
}

// naive error handler for now
func handleError(msg string, err error) error {
	err = fmt.Errorf("%s: %w", msg, err)
	sentry.CaptureException(err)
	return err
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

func contains[T comparable](list []T, item T) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}

// isThreadAlreadyLabeled if the messages already labeled with SRC label
func isThreadAlreadyLabeled(messages []*gmail.Message, srcLabelId string) bool {
	if len(messages) == 0 {
		return true
	}

	// for each message in the thread, check if it has the @src label
	for _, m := range messages {
		if srcmessage.HasLabel(m, srcLabelId) {
			return true
		}
	}

	return false
}

func isMessageBeforeReply(messages []*gmail.Message, messageID string) bool {
	// ensure messages are sorted by ascending date
	srcmessage.SortByDate(messages)

	// skip if the message it doesn't appear before a reply
	for _, m := range messages {
		if srcmessage.IsSent(m) {
			break
		}
		if m.Id == messageID {
			return true
		}
	}
	return false
}

type CloudFunction struct {
	ctx      context.Context
	queries  db.Querier
	srv      *srcmail.Service
	labels   *srclabel.CandidateLabels
	model    ml.Service
	user     db.UserProfile
	settings schema.EmailMessagesSettings
}

func NewCloudFunction(ctx context.Context, payload schema.EmailMessages) (*CloudFunction, error) {
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("error parsing GOOGLE_OAUTH2_CREDENTIALScredentials: %w", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetUserProfileByEmail(ctx, payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile by email: %w", err)
	}

	// 2. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Email:    payload.Email,
		Provider: provider,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user oauth token: %w", err)
	}

	// stop early if user token is marked invalid
	if !userToken.IsValid {
		return nil, fmt.Errorf("user token is not valid: %s", payload.Email)
	}

	// 3. Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := srcmail.NewService(ctx, creds, auth)
	if err != nil {
		return nil, fmt.Errorf("error creating gmail service: %w", err)
	}

	// 4. Get or Create SRC Labels
	labels, err := srv.GetOrCreateCandidateLabels()
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
		if srcmail.IsOAuth2Error(err) {
			log.Printf("error oauth error: %v", err)
			// update the user's oauth token
			err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
				UserID:   userToken.UserID,
				Email:    payload.Email,
				Provider: provider,
				Token:    userToken.Token,
				IsValid:  false,
			})
			if err != nil {
				log.Printf("error updating user oauth token: %v", err)
			} else {
				log.Printf("marked user oauth token as invalid")
			}
		}
		return nil, fmt.Errorf("error getting or creating src labels: %w", err)
	}
	// Create recruiting detector client
	mlServiceBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, mlServiceBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating id token source: %w", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error getting id token: %w", err)
	}

	model := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   mlServiceBaseURL,
		AuthToken: idToken.AccessToken,
	})

	return &CloudFunction{
		ctx:      ctx,
		queries:  queries,
		srv:      srv,
		labels:   labels,
		model:    model,
		user:     user,
		settings: payload.Settings,
	}, nil
}

func handler(ctx context.Context, e event.Event) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       "candidate-gmail-messages",
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %v", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	var msg schema.MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return handleError("error parsing Pub/Sub message", err)
	}

	data := msg.Message.Data
	log.Printf("Event: %s", data)

	var payload schema.EmailMessages
	if err := json.Unmarshal(data, &payload); err != nil {
		return handleError("error parsing email messages", err)
	}

	// validate payload
	// for invalid payloads, we don't want to retry
	if payload.Email == "" || len(payload.Messages) == 0 {
		err = fmt.Errorf("received invalid payload: %v", payload)
		log.Print(err)
		sentry.CaptureException(err)
		return nil
	}

	cf, err := NewCloudFunction(ctx, payload)
	if err != nil {
		return handleError("error creating cloud function", err)
	}

	err = cf.processMessages(payload.Messages)
	if err != nil {
		return handleError("error processing messages", err)
	}

	return nil
}

func (cf *CloudFunction) processMessages(messageIDs []string) error {
	messages := map[string]*gmail.Message{}

	for _, id := range messageIDs {
		// payload isn't included in the list endpoint responses
		message, err := cf.srv.GetMessage(id)
		if err != nil {
			if srcmail.IsNotFound(err) {
				// message was deleted, skip
				log.Printf("skipping message %s was deleted", id)
				continue
			}
			// for now, abort on other errors
			return fmt.Errorf("error getting message %s: %w", id, err)
		}

		// filter out empty messages
		if message.Payload == nil {
			continue
		}

		// filter out messages with the sent or already have a job label
		if contains(message.LabelIds, "SENT") || contains(message.LabelIds, cf.labels.JobsOpportunity.Id) {
			continue
		}

		sender := srcmessage.Sender(message)

		// check if message sender is on the allow list
		allowed, err := cf.srv.IsSenderAllowed(sender)
		if err != nil {
			log.Printf("error checking allow list: %v", err)
			sentry.CaptureException(fmt.Errorf("error checking allow list: %w", err))
		}
		// do not take action on allowed senders
		if allowed {
			log.Printf("allowing message: %s", message.Id)
			continue
		}

		// check if message sender is on the block list
		blocked, err := cf.srv.IsSenderBlocked(sender)
		if err != nil {
			log.Printf("error checking block list: %v", err)
			sentry.CaptureException(fmt.Errorf("error checking block list: %w", err))
		}
		// do not take action on allowed senders
		if blocked {
			err = cf.srv.BlockMessage(message.Id, cf.labels)
			if err != nil {
				log.Printf("error blocking message: %v", err)
				sentry.CaptureException(fmt.Errorf("error blocking message: %w", err))
				continue
			}
			log.Printf("blocked message: %s", message.Id)
			continue
		}
		// get the message thread
		thread, err := cf.srv.GetThread(message.ThreadId, "minimal")
		if err != nil {
			// for now abort on error
			return fmt.Errorf("error getting thread: %w", err)
		}

		// skip thread if we already processed it or the sender has already responded
		if isThreadAlreadyLabeled(thread.Messages, cf.labels.JobsOpportunity.Id) || !isMessageBeforeReply(thread.Messages, message.Id) {
			log.Printf("skipping thread: %s", message.ThreadId)
			continue
		}

		messages[message.Id] = message
	}

	// 1. Check if message ID is a known recruiting outbound message
	// 2. If so, we know it's a recruiting email and we can
	// - skip the prediction step
	// - skip parsing and adding to the database (it's already there)
	// - add the job labels
	knownRecruitingEmailIDs := []string{}
	for _, message := range messages {
		// check if message is a known recruiting outbound message
		internalMessageID := srcmessage.Header(message, "Message-ID")
		if internalMessageID == "" {
			log.Printf("skipping message: %s, no Message-ID header", message.Id)
			continue
		}
		// check if message is a known recruiting outbound message
		// look up by RFC2822 message ID and to_email
		// TODO: Is it better to use the recipient email from the header or the current user?
		recipient := srcmessage.RecipientEmail(message)
		if recipient == "" {
			log.Printf("skipping message: %s, no recipient email", message.Id)
			continue
		}

		_, err := cf.queries.GetRecruiterOutboundMessageByRecipient(cf.ctx, db.GetRecruiterOutboundMessageByRecipientParams{
			ToEmail:           recipient,
			InternalMessageID: internalMessageID,
		})
		// we expect a not found error
		if err == nil {
			// log and add to known messages
			log.Printf("found known recruiting message: %s", message.Id)
			// add to known messages
			knownRecruitingEmailIDs = append(knownRecruitingEmailIDs, message.Id)
			// remove from messages
			delete(messages, message.Id)
			continue
		} else if err != sql.ErrNoRows {
			log.Printf("error looking up known recruiting message: %v", err)
		}
	}

	if len(messages) == 0 && len(knownRecruitingEmailIDs) == 0 {
		return nil
	}

	log.Printf("number of emails to classify: %d", len(messages))

	// Get IDs of new recruiting emails
	recruitingEmailIDs, err := cf.classifyRecruitingEmails(messages)
	// TODO: Support partial failures
	if err != nil {
		return fmt.Errorf("error classifying and parsing messages: %v", err)
	}

	// add known recruiting emails
	recruitingEmailIDs = append(recruitingEmailIDs, knownRecruitingEmailIDs...)

	log.Printf("number of recruiting emails: %d", len(recruitingEmailIDs))

	// Label new recruiting emails
	err = cf.processRecruitingEmails(recruitingEmailIDs)
	if err != nil {
		return fmt.Errorf("error processing recruiting emails: %v", err)
	}

	// Label known recruiting emails (verified oppurtunities)
	err = cf.labelKnownRecruitingEmails(knownRecruitingEmailIDs)
	if err != nil {
		return fmt.Errorf("error processing known recruiting emails (verified oppurtunities): %v", err)
	}

	// Save statistics at end to avoid re-counting
	if !cf.settings.DryRun && len(messages) > 0 {
		err := cf.queries.IncrementUserEmailStat(
			context.Background(),
			db.IncrementUserEmailStatParams{
				UserID:    cf.user.UserID,
				Email:     cf.user.Email,
				StatID:    "emails_processed",
				StatValue: int32(len(messages)),
			},
		)
		if err != nil {
			// print error, but don't abort
			log.Printf("error incrementing user email stat: %v", err)
		}
	}

	return nil
}

func (cf *CloudFunction) classifyRecruitingEmails(messages map[string]*gmail.Message) ([]string, error) {
	// Get IDs of new recruiting emails
	recruitingEmailIDs := []string{}
	if len(messages) == 0 {
		return recruitingEmailIDs, nil
	}

	examples := map[string]*ml.ClassifyRequest{}
	for _, message := range messages {
		examples[message.Id] = &ml.ClassifyRequest{
			From:    srcmessage.Sender(message),
			Subject: srcmessage.Subject(message),
			Body:    srcmessage.Body(message),
		}
	}

	results, err := cf.model.BatchClassify(&ml.BatchClassifyRequest{
		Inputs: examples,
	})
	if err != nil {
		return recruitingEmailIDs, fmt.Errorf("error predicting on examples: %v", err)
	}

	for id, isRecruitingEmail := range results.Results {
		// ignore non-recruiting emails based on classification
		if !isRecruitingEmail {
			continue
		}
		recruitingEmailIDs = append(recruitingEmailIDs, id)
	}

	return recruitingEmailIDs, nil
}

func (cf *CloudFunction) processRecruitingEmails(messageIDs []string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	removeLabels := []string{}
	if cf.user.AutoArchive {
		removeLabels = append(removeLabels, "INBOX", "UNREAD")
	}

	if !cf.settings.DryRun {
		_, err := srcmail.ExecuteWithRetries(func() (interface{}, error) {
			err := cf.srv.Users.Messages.BatchModify(cf.srv.UserID, &gmail.BatchModifyMessagesRequest{
				Ids: messageIDs,
				// Add job opportunity label and parent folder labels
				AddLabelIds:    []string{cf.labels.SRC.Id, cf.labels.Jobs.Id, cf.labels.JobsOpportunity.Id},
				RemoveLabelIds: removeLabels,
			}).Do()
			// hack to make compatible with ExecuteWithRetries requirements
			return nil, err
		})

		if err != nil {
			return fmt.Errorf("error modifying recruiting emails: %v", err)
		}
	}

	return nil
}

func (cf *CloudFunction) labelKnownRecruitingEmails(messageIDs []string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	if !cf.settings.DryRun {
		_, err := srcmail.ExecuteWithRetries(func() (interface{}, error) {
			err := cf.srv.Users.Messages.BatchModify(cf.srv.UserID, &gmail.BatchModifyMessagesRequest{
				Ids: messageIDs,
				// Add verified label
				AddLabelIds:    []string{cf.labels.JobsVerified.Id},
			}).Do()
			// hack to make compatible with ExecuteWithRetries requirements
			return nil, err
		})

		if err != nil {
			return fmt.Errorf("error modifying known (verified) recruiting emails: %v", err)
		}
	}

	return nil
}
