package cloudfunctions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	_ "github.com/lib/pq"

	"github.com/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/libs/gmail"
)

const provider = "google"

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

	log.Printf("Email: %s", email)
	log.Printf("History ID: %d", historyID)

	creds, err := jsonFromEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return err
	}

	// 0, Create SRC client
	connectionURI := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		return err
	}

	queries := client.New(db)

	// 1. Get User from email address
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// 2. Make Request to Fetch Previous History ID
	prevSyncHistory, err := queries.GetUserEmailSyncHistory(ctx, user.ID)
	if err == sql.ErrNoRows {
		log.Printf("no email history found for email: %s", email)
	} else {
		return err
	}

	// 3. Make Request to proactively save new history (If anything goes wrong, then we reset the history ID to the previous one)
	err = queries.UspertUserEmailSyncHistory(ctx, client.UspertUserEmailSyncHistoryParams{
		UserID:    user.ID,
		HistoryID: int64(historyID),
	})
	if err != nil {
		return err
	}

	// if this is the first sync, we are done
	if prevSyncHistory.HistoryID == 0 {
		// TODO: Trigger historic sync on first notification
		return nil
	}

	// 4. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, client.GetUserOAuthTokenParams{
		UserID:   user.ID,
		Provider: provider,
	})
	if err != nil {
		return err
	}

	// 5. Create Gmail Service
	auth := []byte(userToken.Token.RawMessage)
	srv, err := mail.NewGmailService(ctx, creds, auth)
	gmailUser := "me"

	// 6. Make Request to Fetch New Emails from Previous History ID
	messages, err := mail.GetNewEmailsSince(srv, gmailUser, uint64(prevSyncHistory.HistoryID), "INBOX")
	if err != nil {
		return err
	}
	// 7. Stringify Emails
	examples := []string{}
	for _, message := range messages {
		text, err := mail.MessageToString(message)
		if err != nil {
			return err
		}
		examples = append(examples, text)
	}

	log.Printf("New Emails: %d", len(examples))
	// 8. Make Request to Detect Recruiting Emails
	// 9. Get or Create SRC Label
	// 10. Batch Modify Emails

	return nil
}
