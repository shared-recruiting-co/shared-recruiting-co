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
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const (
	provider = "google"
)

var (
	// global variable to share across functions...simplest approach for now
	examplesCollectorSrv   *srcmail.Service
	collectedExampleLabels = []string{"INBOX", "UNREAD"}
)

func init() {
	functions.CloudEvent("Handler", handler)
}

type EmailHistory struct {
	Email     string `json:"email"`
	HistoryID int64  `json:"historyId"`
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

// naive error handler for now
func handleError(msg string, err error) error {
	err = fmt.Errorf("%s: %w", msg, err)
	sentry.CaptureException(err)
	return err
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

	email := emailPushNotification.Email
	historyID := emailPushNotification.HistoryID

	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return handleError("error parsing GOOGLE_OAUTH2_CREDENTIALS", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetUserProfileByEmail(ctx, email)
	if err != nil {
		return handleError("error getting user profile by email", err)
	}
	if user.AutoContribute {
		auth, err := jsonFromEnv("EXAMPLES_GMAIL_OAUTH_TOKEN")
		if err != nil {
			return handleError("error parsing EXAMPLES_GMAIL_OAUTH_TOKEN", err)
		}
		examplesCollectorSrv, err = srcmail.NewService(ctx, creds, auth)
		if err != nil {
			return handleError("error creating examples collector service", err)
		}
	}

	// 2. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Provider: provider,
	})
	if err != nil {
		return handleError("error getting user oauth token", err)
	}

	// stop early if user token is marked invalid
	if !userToken.IsValid {
		log.Printf("user token is not valid: %s", email)
		return nil
	}

	// 3. Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := srcmail.NewService(ctx, creds, auth)
	if err != nil {
		return handleError("error creating gmail service", err)
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
		return handleError("error getting or creating SRC labels", err)
	}
	// Create recruiting detector client
	classifierBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, classifierBaseURL)
	if err != nil {
		return handleError("error creating id token source", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return handleError("error getting id token", err)
	}

	classifier := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   classifierBaseURL,
		AuthToken: idToken.AccessToken,
	})

	// 5. Make Request to get previous history and proactively save new history (If anything goes wrong, then we reset the history ID to the previous one)
	// Make Request to Fetch Previous History ID
	prevSyncHistory, err := queries.GetUserEmailSyncHistory(ctx, user.UserID)
	// On first notification, trigger a full sync in the background
	if err == sql.ErrNoRows {
		log.Printf("no previous sync history found, triggering full sync in background")
		// let's sync one year of emails for now
		startDate := time.Now().AddDate(-1, 0, 0)
		err = triggerBackgroundfFullEmailSync(ctx, email, startDate)
		if err != nil {
			return handleError("error triggering background full email sync", err)
		}

		// save the current history ID
		err = queries.UpsertUserEmailSyncHistory(ctx, db.UpsertUserEmailSyncHistoryParams{
			UserID:    user.UserID,
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

	err = queries.UpsertUserEmailSyncHistory(ctx, db.UpsertUserEmailSyncHistoryParams{
		UserID:    user.UserID,
		HistoryID: int64(historyID),
		SyncedAt:  time.Now(),
	})
	if err != nil {
		return handleError("error upserting user email sync history", err)
	}

	// on any errors after this, we want to reset the history ID to the previous one
	revertSynctHistory := func() {
		err := queries.UpsertUserEmailSyncHistory(ctx, db.UpsertUserEmailSyncHistoryParams{
			UserID:    user.UserID,
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
			revertSynctHistory()
		}
	}()

	// 7. Sync new emails
	err = syncNewEmails(
		user,
		srv,
		queries,
		classifier,
		prevSyncHistory,
		labels,
	)
	if err != nil {
		revertSynctHistory()
		return handleError("error syncing new emails", err)
	}

	log.Printf("done.")
	return nil
}
