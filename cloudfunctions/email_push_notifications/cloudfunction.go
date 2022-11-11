package cloudfunctions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"

	"github.com/shared-recruiting-co/libs/gmail"
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
	HistoryID int64  `json:"historyId"`
}

type EmailHistory struct {
	Email     string `json:"email"`
	HistoryID int64  `json:"historyId"`
}

func JsonFromEnv(env string) ([]byte, error) {
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

	// TODO:
	// 1. Make Request to Fetch Previous History ID
	// 2. Make Request to Save New History ID (If anything goes wrong, then we reset the history ID to the previous one)
	// 3. Fetch the user's access token
	// 4. Create a Gmail service
	creds, err := JsonFromEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return err
	}
	auth, err := JsonFromEnv("GOOGLE_AUTH_TOKEN")
	if err != nil {
		return err
	}
	srv, err := gmail.NewGmailService(ctx, creds, auth)
	// 5. Make Request to Fetch New Emails from Previous History ID
	messages, err := gmail.GetNewEmailsSince(srv, historyID, "INBOX")
	if err != nil {
		return err
	}
	// 6. Stringify Emails
	examples := []string{}
	for _, message := range messages {
		text, err := gmail.MessageToString(message)
		if err != nil {
			return err
		}
		examples = append(examples, text)
	}
	// 7. Make Request to Detect Recruiting Emails
	// 8. Get or Create SRC Label
	// 9. Batch Modify Emails

	return nil
}
