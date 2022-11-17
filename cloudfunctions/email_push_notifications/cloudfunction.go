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

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
)

const (
	provider                = "google"
	SRC_Label               = "@SRC"
	SRC_JobOpportunityLabel = "@SRC/Job Opportunity"
	SRC_Color               = "#16a765"
	White                   = "#ffffff"
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

	// prepare queries
	queries, err := client.Prepare(ctx, db)
	if err != nil {
		return err
	}

	// 1. Get User from email address
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// 2. Make Request to Fetch Previous History ID
	prevSyncHistory, err := queries.GetUserEmailSyncHistory(ctx, user.ID)
	if err == sql.ErrNoRows {
		log.Printf("no email history found for email: %s", email)
	} else if err != nil {
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
	gmailSrv, err := mail.NewGmailService(ctx, creds, auth)
	gmailUser := "me"

	// 6. Make Request to Fetch New Emails from Previous History ID
	messages, err := mail.GetNewEmailsSince(gmailSrv, gmailUser, uint64(prevSyncHistory.HistoryID), "INBOX")
	if err != nil {
		return err
	}

	// 7. Stringify Emails
	examples := map[string]string{}
	for _, message := range messages {
		text, err := mail.MessageToString(message)
		if err != nil {
			return err
		}
		examples[message.Id] = text
	}

	log.Printf("number of new emails: %d", len(examples))

	if len(examples) == 0 {
		return nil
	}

	// 8. Create recruiting detector client
	classifier := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL: os.Getenv("CLASSIFIER_URL"),
		ApiKey:  os.Getenv("CLASSIFIER_API_KEY"),
	})

	// 9. Batch predict on new emails
	results, err := classifier.PredictBatch(examples)
	if err != nil {
		return err
	}

	// 10. Get IDs of new recruiting emails
	recruitingEmailIDs := []string{}
	for id, result := range results {
		if !result {
			continue
		}
		recruitingEmailIDs = append(recruitingEmailIDs, id)
	}

	// no new recruiting emails, return early
	if len(recruitingEmailIDs) == 0 {
		return nil
	}

	// 11. Get or Create SRC Label
	srcLabel, err := mail.GetOrCreateLabel(gmailSrv, gmailUser, SRC_Label, SRC_Color, White)
	srcJobOpportunityLabel, err := mail.GetOrCreateLabel(gmailSrv, gmailUser, SRC_JobOpportunityLabel, SRC_Color, White)

	// 12. Take action! (Batch Modify Emails)
	err = gmailSrv.Users.Messages.BatchModify(gmailUser, &gmail.BatchModifyMessagesRequest{
		Ids: recruitingEmailIDs,
		// Add SRC Label
		AddLabelIds: []string{srcLabel.Id, srcJobOpportunityLabel.Id},
		// In future,
		// - mark as read
		// - archive
		// - create response
		// RemoveLabelIds: []string{"UNREAD"},
	}).Do()

	if err != nil {
		return err
	}

	return nil
}
