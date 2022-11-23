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
	provider                     = "google"
	SRC_Label                    = "@SRC"
	SRC_JobOpportunityLabel      = "@SRC/Job Opportunity"
	SRC_Color                    = "#ff7537"
	SRC_JobOpportunityLabelColor = "#16a765"
	White                        = "#ffffff"
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

	creds, err := jsonFromEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		return fmt.Errorf("error fetching google app credentials: %v", err)
	}

	// 0, Create SRC client
	connectionURI := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()
	// use a max of 2 connections
	db.SetMaxOpenConns(2)

	// prepare queries
	queries, err := client.Prepare(ctx, db)
	if err != nil {
		return fmt.Errorf("error preparing db queries: %v", err)
	}

	// 1. Get User from email address
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error getting user from email: %v", err)
	}

	// 2. Make Request to Fetch Previous History ID
	prevSyncHistory, err := queries.GetUserEmailSyncHistory(ctx, user.ID)
	if err == sql.ErrNoRows {
		log.Printf("no email history found for email: %s", email)
	} else if err != nil {
		return fmt.Errorf("error getting user email sync history: %v", err)
	}

	// 3. Make Request to proactively save new history (If anything goes wrong, then we reset the history ID to the previous one)
	err = queries.UspertUserEmailSyncHistory(ctx, client.UspertUserEmailSyncHistoryParams{
		UserID:    user.ID,
		HistoryID: int64(historyID),
	})
	if err != nil {
		return fmt.Errorf("error upserting email sync history: %v", err)
	}

	// 4. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, client.GetUserOAuthTokenParams{
		UserID:   user.ID,
		Provider: provider,
	})
	if err != nil {
		return fmt.Errorf("error getting user oauth token: %v", err)
	}

	// 5. Create Gmail Service
	auth := []byte(userToken.Token.RawMessage)
	gmailSrv, err := mail.NewGmailService(ctx, creds, auth)
	gmailUser := "me"

	// 6. Get or Create SRC Label
	srcLabel, err := mail.GetOrCreateLabel(gmailSrv, gmailUser, SRC_Label, SRC_Color, White)
	if err != nil {
		return fmt.Errorf("error getting or creating the SRC label: %v", err)
	}
	srcJobOpportunityLabel, err := mail.GetOrCreateLabel(gmailSrv, gmailUser, SRC_JobOpportunityLabel, SRC_JobOpportunityLabelColor, White)
	if err != nil {
		return fmt.Errorf("error getting or creating the SRC job opportunity label: %v", err)
	}

	// 7. Create recruiting detector client
	classifier := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL: os.Getenv("CLASSIFIER_URL"),
		ApiKey:  os.Getenv("CLASSIFIER_API_KEY"),
	})

	var messages []*gmail.Message
	pageToken := ""

	// batch process messages to reduce memory usage
	for {

		// 8. Make Request to Fetch New Emails from Previous History ID
		// get next set of messages
		// if this is the first notification, trigger a one-time sync
		if prevSyncHistory.HistoryID == 0 {
			messages, pageToken, err = GetEmailsSinceLastYear(gmailSrv, gmailUser, pageToken)
		} else {
			messages, pageToken, err = GetNewEmailsSince(gmailSrv, gmailUser, uint64(prevSyncHistory.HistoryID), "INBOX", pageToken)
		}

		// for now, abort on error
		if err != nil {
			return fmt.Errorf("error fetching emails: %v", err)
		}

		// process messages
		examples := map[string]string{}
		for _, message := range messages {
			// payload isn't included in the list endpoint responses
			message, err := gmailSrv.Users.Messages.Get(gmailUser, message.Id).Do()

			// for now, abort on error
			if err != nil {
				return fmt.Errorf("error getting message: %v", err)
			}

			if message.Payload == nil {
				continue
			}
			text, err := mail.MessageToString(message)
			examples[message.Id] = text
		}

		log.Printf("number of emails to classify: %d", len(examples))

		if len(examples) == 0 {
			break
		}

		// 9. Batch predict on new emails
		results, err := classifier.PredictBatch(examples)
		if err != nil {
			return fmt.Errorf("error predicting on examples: %v", err)
		}

		// 10. Get IDs of new recruiting emails
		recruitingEmailIDs := []string{}
		for id, result := range results {
			if !result {
				continue
			}
			recruitingEmailIDs = append(recruitingEmailIDs, id)
		}

		log.Printf("number of recruiting emails: %d", len(recruitingEmailIDs))

		// 11. Take action on recruiting emails
		if len(recruitingEmailIDs) > 0 {
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

			// for now, abort on error
			if err != nil {
				return fmt.Errorf("error modifying recruiting emails: %v", err)
			}
		}

		if pageToken == "" {
			break
		}
	}

	log.Printf("done.")
	return nil
}
