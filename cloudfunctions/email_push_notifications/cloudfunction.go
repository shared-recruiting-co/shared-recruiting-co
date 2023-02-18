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

	"cloud.google.com/go/pubsub"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/getsentry/sentry-go"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const (
	provider = "google"
)

func init() {
	functions.CloudEvent("Handler", handler)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

// naive error handler for now
func handleError(msg string, err error) error {
	err = fmt.Errorf("%s: %w", msg, err)
	sentry.CaptureException(err)
	return err
}

type PubSubTopics struct {
	CandidateGmailMessages *pubsub.Topic
}

type CloudFunction struct {
	ctx     context.Context
	queries db.Querier
	srv     *srcmail.Service
	labels  *srclabel.Labels
	user    db.UserProfile
	payload schema.EmailPushNotification
	topics  *PubSubTopics
}

func NewCloudFunction(ctx context.Context, payload schema.EmailPushNotification) (*CloudFunction, error) {
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

	projectID := os.Getenv("GCP_PROJECT_ID")
	if err != nil {
		return nil, err
	}
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		// TODO: Handle error.
		log.Printf("failed to create pubsub client: %v", err)
	}
	candidateGmailMessagesTopic := os.Getenv("CANDIDATE_GMAIL_MESSAGES_TOPIC")
	topic := client.Topic(candidateGmailMessagesTopic)

	topics := &PubSubTopics{
		CandidateGmailMessages: topic,
	}

	return &CloudFunction{
		ctx:     ctx,
		queries: queries,
		srv:     srv,
		labels:  labels,
		user:    user,
		payload: payload,
		topics:  topics,
	}, nil
}

// handler consumes a CloudEvent message and extracts the Pub/Sub message.
func handler(ctx context.Context, e event.Event) error {
	var msg schema.MessagePublishedData

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       "email-push-notifications",
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %v", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	if err := e.DataAs(&msg); err != nil {
		return handleError("error parsing Pub/Sub message", err)
	}

	data := msg.Message.Data
	log.Printf("Event: %s", data)

	var emailPushNotification schema.EmailPushNotification
	if err := json.Unmarshal(data, &emailPushNotification); err != nil {
		return handleError("error parsing email push notification", err)
	}

	historyID := emailPushNotification.HistoryID

	cf, err := NewCloudFunction(ctx, emailPushNotification)
	if err != nil {
		return handleError("error initializing cloud function", err)
	}

	// 5. Make Request to get previous history and proactively save new history (If anything goes wrong, then we reset the history ID to the previous one)
	// Make Request to Fetch Previous History ID
	prevSyncHistory, err := cf.queries.GetUserEmailSyncHistory(ctx, db.GetUserEmailSyncHistoryParams{
		UserID:    cf.user.UserID,
		InboxType: db.InboxTypeCandidate,
		Email:     cf.payload.Email,
	})
	// On first notification, trigger a full sync in the background
	if err == sql.ErrNoRows {
		log.Printf("no previous sync history found, triggering full sync in background")
		// let's sync one year of emails for now
		startDate := time.Now().AddDate(-1, 0, 0)
		err = cf.triggerBackgroundfFullEmailSync(startDate)
		if err != nil {
			return handleError("error triggering background full email sync", err)
		}

		// save the current history ID
		err = cf.queries.UpsertUserEmailSyncHistory(ctx, db.UpsertUserEmailSyncHistoryParams{
			UserID:    cf.user.UserID,
			InboxType: db.InboxTypeCandidate,
			Email:     cf.payload.Email,
			HistoryID: int64(historyID),
			SyncedAt:  time.Now(),
		})
		if err != nil {
			return handleError("error upserting user email sync history", err)
		}
		// success
		log.Printf("done.")
		return nil

	} else if err != nil {
		return handleError("error getting user email sync history", err)
	}

	err = cf.queries.UpsertUserEmailSyncHistory(ctx, db.UpsertUserEmailSyncHistoryParams{
		UserID:    cf.user.UserID,
		InboxType: db.InboxTypeCandidate,
		Email:     cf.payload.Email,
		HistoryID: int64(historyID),
		SyncedAt:  time.Now(),
	})
	if err != nil {
		return handleError("error upserting user email sync history", err)
	}

	// on any errors after this, we want to reset the history ID to the previous one
	revertSyncHistory := func() {
		err := cf.queries.UpsertUserEmailSyncHistory(ctx, db.UpsertUserEmailSyncHistoryParams{
			UserID:    cf.user.UserID,
			InboxType: db.InboxTypeCandidate,
			Email:     cf.payload.Email,
			HistoryID: prevSyncHistory.HistoryID,
			SyncedAt:  prevSyncHistory.SyncedAt,
		})
		if err != nil {
			log.Printf("error reverting email sync history: %v", err)
		}
	}

	// revert sync history if we hit an unexpected error past this point
	// Note: deferred functions are called in LIFO order, so this will be called before the defer db.Close()
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
			log.Println("reverting sync history")
			revertSyncHistory()
		}
	}()

	// 7. Sync new emails
	err = cf.syncHistory(prevSyncHistory)
	if err != nil {
		revertSyncHistory()
		return handleError("error syncing new emails", err)
	}

	log.Printf("done.")
	return nil
}
