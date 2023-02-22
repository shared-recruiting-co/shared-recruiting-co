package cloudfunctions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/jaytaylor/html2text"

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

type CloudFunction struct {
	ctx      context.Context
	queries  db.Querier
	srv      *srcmail.Service
	labels   *srclabel.RecruiterLabels
	model    ml.Service
	user     db.GetRecruiterByEmailRow
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
	user, err := queries.GetRecruiterByEmail(ctx, payload.Email)
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

	// Get the SRC labels
	labels, err := srv.GetOrCreateRecruiterLabels()
	if err != nil {
		return nil, fmt.Errorf("error getting recruiter labels: %w", err)
	}

	// Create recruiting detector client
	classifierBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, classifierBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating id token source: %w", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error getting id token: %w", err)
	}

	model := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   classifierBaseURL,
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
		ServerName:       os.Getenv("FUNCTION_NAME"),
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

	for _, m := range payload.Messages {
		err = cf.processMessage(m)
		if err != nil {
			// for now abort on first error
			return handleError("error creating cloud function", err)
		}
	}

	return nil
}

// processMessage processes a single message
//
// Logic
// 1. Given message, get thread
// 2. Check if thread is already labeled if so skip
// 3. Get the first message in the thread
//
// 4. Check if the first message in the thread exists in the database
// 5. If it does, skip
// 6. If it doesn't check for a matching template
// 7. If there is a matching template, save the message to the database and label
// 8. If not, check if the message is a recruiting email (model.Classify)
// 9. If not, skip
// 10. If it is, save the message as a template, save the message to the database and label
func (cf *CloudFunction) processMessage(id string) error {
	// 1. Get Message
	msg, err := cf.srv.GetMessage(id)
	if err != nil {
		return fmt.Errorf("error getting message: %w", err)
	}

	// 2. Get Thread
	thread, err := cf.srv.GetThread(msg.ThreadId, "minimal")
	if err != nil {
		return fmt.Errorf("error getting thread: %w", err)
	}

	// 3. Check if thread is already labeled or if the earliest message is not from the user
	if !isThreadStartedByUser(thread.Messages) || isThreadLabeled(thread.Messages, cf.labels.RecruitingOutbound.Id) {
		log.Printf("thread %s already labeled or not started by user, skipping", thread.Id)
		return nil
	}

	// 4. Get the first message in the thread
	// Two options:
	// a. There are multiple messages in the thread
	// b. There is only one message in the thread. In this case, we search for previous emails to the recipient
	//    and use the first one as the first message in the thread
	var firstMsg *gmail.Message
	if len(thread.Messages) > 1 {
		srcmessage.SortByDate(thread.Messages)
		firstMsg = thread.Messages[0]
	} else {
		// TODO: search for previous emails to the recipient within 3 months
		firstMsg = msg
	}

	// 5. Check if the first message in the thread exists in the database
	existing, err := cf.queries.GetRecruiterOutboundMessage(cf.ctx, db.GetRecruiterOutboundMessageParams{
		RecruiterID: cf.user.UserID,
		MessageID:   firstMsg.Id,
	})

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking if first message exists in the database: %w", err)
	}
	if existing.MessageID != "" {
		log.Printf("first message %s already exists in the database, skipping", existing.MessageID)
		return nil
	}

	if firstMsg.Id != msg.Id {
		firstMsg, err = cf.srv.GetMessage(firstMsg.Id)
		if err != nil {
			return fmt.Errorf("error getting first message: %w", err)
		}
	}
	// extract first message subject and body text
	subject := srcmessage.Subject(firstMsg)
	body := srcmessage.Body(firstMsg)
	normalizedBody := normalizeBody(body)
	normalizedContent := fmt.Sprintf("%s %s", subject, normalizedBody)

	// 6. If it doesn't check for a matching template
	templates, err := cf.queries.ListSimilarRecruiterOutboundTemplates(cf.ctx, db.ListSimilarRecruiterOutboundTemplatesParams{
		UserID: cf.user.UserID,
		Input:  normalizedContent,
	})
	if err != nil {
		return fmt.Errorf("error getting similar templates: %w", err)
	}

	// matching template
	// 7. If there is a matching template, save the message to the database and label
	if len(templates) > 0 {
		// for now, print all matching templates and match to the first one
		log.Printf("found %d matching templates", len(templates))
		for _, t := range templates {
			log.Printf("template: %s / similarity %0.2f", t.TemplateID, t.Similarity)
		}

		// match to the first template
		template := templates[0]

		err = cf.queries.InsertRecruiterOutboundMessage(cf.ctx, db.InsertRecruiterOutboundMessageParams{
			RecruiterID: cf.user.UserID,
			MessageID:   firstMsg.Id,
			TemplateID: uuid.NullUUID{
				UUID:  template.TemplateID,
				Valid: true,
			},
			InternalMessageID: srcmessage.HostMessageID(firstMsg),
			FromEmail:         srcmessage.SenderEmail(firstMsg),
			ToEmail:           srcmessage.RecipientEmail(firstMsg),
			SentAt:            srcmessage.CreatedAt(firstMsg),
		})
		if err != nil {
			return fmt.Errorf("error inserting message: %w", err)
		}

		// label thread
		_, err = srcmail.ExecuteWithRetries(func() (interface{}, error) {
			return cf.srv.Users.Threads.Modify(cf.srv.UserID, thread.Id, &gmail.ModifyThreadRequest{
				// Add job opportunity label and parent folder labels
				AddLabelIds: []string{cf.labels.SRC.Id, cf.labels.Recruiting.Id, cf.labels.RecruitingOutbound.Id},
			}).Do()
		})
		if err != nil {
			return fmt.Errorf("error labeling thread: %w", err)
		}

		// done
		return nil
	}
	log.Printf("no matching template found for message %s", msg.Id)

	// 8. If not, check if the message is a recruiting email (model.Classify)
	classification, err := cf.model.Classify(&ml.EmailInput{
		From:    srcmessage.Sender(firstMsg),
		Subject: subject,
		// Use raw body
		Body: body,
	})
	if err != nil {
		return fmt.Errorf("error classifying email: %w", err)
	}
	// 9. If not, skip
	if !classification.Result {
		log.Printf("first message %s is not a recruiting email, skipping", firstMsg.Id)
		return nil
	}

	// 10. If it is, save the message as a template, save the message to the database and label
	metadata := map[string]interface{}{
		"MessageID": firstMsg.Id,
		"SentAt":    srcmessage.CreatedAt(firstMsg),
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error marshaling template metadata: %w", err)
	}
	template, err := cf.queries.InsertRecruiterOutboundTemplate(cf.ctx, db.InsertRecruiterOutboundTemplateParams{
		RecruiterID:       cf.user.UserID,
		Subject:           subject,
		Body:              body,
		NormalizedContent: normalizedContent,
		Metadata:          metadataJSON,
	})
	if err != nil {
		return fmt.Errorf("error inserting template: %w", err)
	}

	err = cf.queries.InsertRecruiterOutboundMessage(cf.ctx, db.InsertRecruiterOutboundMessageParams{
		RecruiterID: cf.user.UserID,
		MessageID:   firstMsg.Id,
		TemplateID: uuid.NullUUID{
			UUID:  template.TemplateID,
			Valid: true,
		},
		InternalMessageID: srcmessage.HostMessageID(firstMsg),
		FromEmail:         srcmessage.SenderEmail(firstMsg),
		ToEmail:           srcmessage.RecipientEmail(firstMsg),
		SentAt:            srcmessage.CreatedAt(firstMsg),
	})
	if err != nil {
		return fmt.Errorf("error inserting message: %w", err)
	}

	// label thread
	_, err = srcmail.ExecuteWithRetries(func() (interface{}, error) {
		return cf.srv.Users.Threads.Modify(cf.srv.UserID, thread.Id, &gmail.ModifyThreadRequest{
			// Add job opportunity label and parent folder labels
			AddLabelIds: []string{cf.labels.SRC.Id, cf.labels.Recruiting.Id, cf.labels.RecruitingOutbound.Id},
		}).Do()
	})
	if err != nil {
		return fmt.Errorf("error labeling thread: %w", err)
	}

	return nil
}

// isThreadLabeled if the messages already labeled with SRC label
func isThreadLabeled(messages []*gmail.Message, labelID string) bool {
	if len(messages) == 0 {
		return true
	}

	// for each message in the thread, check if it has the @src label
	for _, m := range messages {
		if srcmessage.HasLabel(m, labelID) {
			return true
		}
	}

	return false
}

func isThreadStartedByUser(messages []*gmail.Message) bool {
	if len(messages) == 0 {
		return false
	}
	// get the first message in the thread
	srcmessage.SortByDate(messages)
	firstMsg := messages[0]
	return srcmessage.IsSent(firstMsg)
}

func normalizeBody(body string) string {
	// convert to text if it's HTML
	mime := http.DetectContentType([]byte(body))
	// mime will be "text/html; charset=utf-8"
	if strings.HasPrefix(mime, "text/html") {
		text, err := html2text.FromString(body,
			html2text.Options{
				PrettyTables: false,
				OmitLinks:    true,
				TextOnly:     true,
			})
		if err != nil {
			// log
			log.Printf("error converting html to text: %v", err)
		} else {
			body = text
		}
	}
	// remove html tags
	body = html.UnescapeString(body)

	// remove extra spaces
	body = strings.ReplaceAll(body, "\n", " ")
	body = strings.ReplaceAll(body, "\t", " ")
	body = strings.ReplaceAll(body, "  ", " ")

	// strip links because they are often customized per recipient
	re := regexp.MustCompile(`https?://\S+`)
	body = re.ReplaceAllString(body, "")

	return body
}
