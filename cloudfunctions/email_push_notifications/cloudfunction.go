package cloudfunctions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
)

const (
	provider = "google"
)

func init() {
	functions.CloudEvent("EmailPushNotificationHandler", emailPushNotificationHandler)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type EmailPushNotification struct {
	Email     string `json:"emailAddress"`
	HistoryID uint64 `json:"historyId"`
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

// emailPushNotificationHandler consumes a CloudEvent message and extracts the Pub/Sub message.
func emailPushNotificationHandler(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	data := msg.Message.Data
	log.Printf("Event: %s", data)

	var emailPushNotification EmailPushNotification
	if err := json.Unmarshal(data, &emailPushNotification); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	email := emailPushNotification.Email
	historyID := emailPushNotification.HistoryID

	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return fmt.Errorf("error fetching google app credentials: %v", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := client.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetUserProfileByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error getting user profile by email: %v", err)
	}

	// 2. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, client.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Provider: provider,
	})
	if err != nil {
		return fmt.Errorf("error getting user oauth token: %v", err)
	}

	// stop early if user token is marked invalid
	if !userToken.IsValid {
		log.Printf("user token is not valid: %s", email)
		return nil
	}

	// 3. Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := mail.NewService(ctx, creds, auth)
	if err != nil {
		return fmt.Errorf("error creating gmail service: %v", err)
	}

	// 4. Get or Create SRC Labels
	_, err = srv.GetOrCreateSRCLabel()
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
		oauth2Err := &oauth2.RetrieveError{}
		if errors.As(err, &oauth2Err) {
			log.Printf("error oauth error: %v", oauth2Err)
			// update the user's oauth token
			err = queries.UpsertUserOAuthToken(ctx, client.UpsertUserOAuthTokenParams{
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
		return fmt.Errorf("error getting or creating the SRC label: %v", err)
	}
	srcJobOpportunityLabel, err := srv.GetOrCreateSRCJobOpportunityLabel()
	if err != nil {
		return fmt.Errorf("error getting or creating the SRC job opportunity label: %v", err)
	}

	// Create recruiting detector client
	classifierBaseURL := os.Getenv("CLASSIFIER_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, classifierBaseURL)
	if err != nil {
		return fmt.Errorf("error creating id token source: %v", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return fmt.Errorf("error getting id token: %v", err)
	}

	classifier := NewClassifierClient(ctx, ClassifierClientArgs{
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
			return fmt.Errorf("error triggering background full email sync: %v", err)
		}

		// save the current history ID
		err = queries.UpsertUserEmailSyncHistory(ctx, client.UpsertUserEmailSyncHistoryParams{
			UserID:              user.UserID,
			HistoryID:           int64(historyID),
			SyncedAt:            time.Now(),
			ExamplesCollectedAt: prevSyncHistory.ExamplesCollectedAt,
		})
		if err != nil {
			return fmt.Errorf("error upserting email sync history: %v", err)
		}
		// success
		log.Printf("done.")
		return nil

	} else if err != nil {
		return fmt.Errorf("error getting user email sync history: %v", err)
	}

	err = queries.UpsertUserEmailSyncHistory(ctx, client.UpsertUserEmailSyncHistoryParams{
		UserID:              user.UserID,
		HistoryID:           int64(historyID),
		SyncedAt:            time.Now(),
		ExamplesCollectedAt: prevSyncHistory.ExamplesCollectedAt,
	})
	if err != nil {
		return fmt.Errorf("error upserting email sync history: %v", err)
	}

	// on any errors after this, we want to reset the history ID to the previous one
	revertSynctHistory := func() {
		err := queries.UpsertUserEmailSyncHistory(ctx, client.UpsertUserEmailSyncHistoryParams{
			UserID:              user.UserID,
			HistoryID:           prevSyncHistory.HistoryID,
			SyncedAt:            prevSyncHistory.SyncedAt,
			ExamplesCollectedAt: prevSyncHistory.ExamplesCollectedAt,
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
	err = syncNewEmails(email, srv, classifier, prevSyncHistory, srcJobOpportunityLabel.Id)
	if err != nil {
		revertSynctHistory()
		return fmt.Errorf("error processing messages: %v", err)
	}

	log.Printf("done.")
	return nil
}
