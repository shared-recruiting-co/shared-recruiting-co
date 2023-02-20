package cloudfunctions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/getsentry/sentry-go"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const (
	provider = "google"
)

func init() {
	functions.HTTP("Handler", handler)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

// generic error handler
func handleError(w http.ResponseWriter, msg string, err error) {
	err = fmt.Errorf("%s: %w", msg, err)
	log.Print(err)
	sentry.CaptureException(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type EmailSyncRequest struct {
	Email string `json:"email"`
	// StartDate is the date to start syncing from (inclusive)
	StartDate time.Time `json:"start_date"`
	// EndDate is optional
	// if not provided, it will default to the current time
	EndDate  time.Time                    `json:"end_date"`
	Settings schema.EmailMessagesSettings `json:"settings"`
}

type CloudFunction struct {
	ctx     context.Context
	queries db.Querier
	srv     *srcmail.Service
	labels  *srclabel.Labels
	user    db.Recruiter
	request EmailSyncRequest
	topic   *pubsub.Topic
}

func NewCloudFunction(ctx context.Context, payload EmailSyncRequest) (*CloudFunction, error) {
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("error parsing GOOGLE_OAUTH2_CREDENTIALS: %w", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetRecruiterByEmail(ctx, payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting recruiter by email: %w", err)
	}

	// 2. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Provider: provider,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user oauth token: %w", err)
	}

	//
	if !userToken.IsValid {
		return nil, fmt.Errorf("user token is not valid: %s", userToken.UserID)
	}

	// 3. Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := srcmail.NewService(ctx, creds, auth)
	if err != nil {
		return nil, fmt.Errorf("error creating gmail service: %w", err)
	}

	// 4. Get or Create SRC Labels
	labels, err := srv.GetOrCreateSRCLabels()
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
		if srcmail.IsOAuth2Error(err) {
			log.Printf("error oauth error: %v", err)
			// update the user's oauth token
			err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
				UserID:   userToken.UserID,
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

	// Create messages topic
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT_ID is not set")
	}
	topicName := os.Getenv("GMAIL_MESSAGES_TOPIC")
	if topicName == "" {
		return nil, fmt.Errorf("GMAIL_MESSAGES_TOPIC is not set")
	}
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("error creating pubsub client: %w", err)
	}
	topic := client.Topic(topicName)
	if topic == nil {
		return nil, fmt.Errorf("error getting pubsub topic: %w", err)
	}

	return &CloudFunction{
		ctx:     ctx,
		queries: queries,
		srv:     srv,
		labels:  labels,
		user: db.Recruiter{
			UserID:    user.UserID,
			CompanyID: user.CompanyID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		request: payload,
		topic:   topic,
	}, nil
}

// handler is triggers a sync up to the specified date for the given email
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       os.Getenv("FUNCTION_NAME"),
	})
	if err != nil {
		log.Printf("sentry.Init: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	// Get the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, "error reading request body", err)
		return
	}
	var data EmailSyncRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		handleError(w, "error unmarshalling request body", err)
		return
	}
	log.Println("email sync request")

	cf, err := NewCloudFunction(ctx, data)
	if err != nil {
		handleError(w, "error creating cloud function", err)
		return
	}
	defer cf.topic.Stop()

	err = cf.Sync()
	if err != nil {
		handleError(w, "error syncing", err)
		return
	}

	log.Printf("done.")
}

func (cf *CloudFunction) Sync() error {
	var err error
	var threads []*gmail.Thread
	var results []*pubsub.PublishResult
	pageToken := ""

	// batch process messages to reduce memory usage
	for {
		// get next set of messages
		// if this is the first notification, trigger a one-time sync
		threads, pageToken, err = fetchThreadsSinceDate(cf.srv, cf.request.StartDate, cf.request.EndDate, pageToken)

		// for now, abort on error
		if err != nil {
			return fmt.Errorf("error fetching emails: %w", err)
		}

		// get the messages for each thread
		var messages []*gmail.Message
		for _, t := range threads {
			thread, err := cf.srv.GetThread(t.Id, "minimal")
			if err != nil {
				// for now abort on error
				return fmt.Errorf("error fetching thread: %w", err)
			}

			// check if we already processed this thread
			if skipThread(thread.Messages, cf.labels.JobsOpportunity.Id) {
				continue
			}
			// (for now) we only want to check the first message in a thread
			srcmessage.SortByDate(messages)
			// save for processing
			messages = append(messages, thread.Messages[0])
		}

		if len(messages) == 0 {
			log.Printf("no messages to process")
		} else {
			// if there are messages to process, push them
			result, err := cf.PublishMessages(messages)
			if err != nil {
				return fmt.Errorf("error publishing messages: %w", err)
			}

			results = append(results, result)
		}

		if pageToken == "" {
			break
		}
	}

	// wait for all messages to be processed
	for _, result := range results {
		_, err := result.Get(cf.ctx)
		if err != nil {
			// log but do not abort
			log.Printf("error getting publish result: %v", err)
		}
	}

	return nil
}

func (cf *CloudFunction) PublishMessages(messages []*gmail.Message) (*pubsub.PublishResult, error) {
	emailMessages := schema.EmailMessages{
		Email:    cf.user.Email,
		Messages: make([]string, len(messages)),
		Settings: cf.request.Settings,
	}
	for i, message := range messages {
		emailMessages.Messages[i] = message.Id
	}
	rawMessage, err := json.Marshal(emailMessages)
	if err != nil {
		return nil, fmt.Errorf("error marshalling email messages: %w", err)
	}

	// publish message
	result := cf.topic.Publish(cf.ctx, &pubsub.Message{
		Data: rawMessage,
	})

	return result, nil
}
