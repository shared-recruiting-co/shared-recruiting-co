package cloudfunctions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

type FullEmailSyncRequest struct {
	Email     string    `json:"email"`
	StartDate time.Time `json:"start_date"`
}

func publishMessages(email string, messages []*gmail.Message) {
	if len(messages) == 0 {
		return
	}

	ctx := context.Background()
	// push message to be processed
	emailMessages := schema.EmailMessages{
		Email:    email,
		Messages: make([]string, len(messages)),
		Settings: schema.EmailMessagesSettings{},
	}
	for i, message := range messages {
		emailMessages.Messages[i] = message.Id
	}
	projectID := os.Getenv("GCP_PROJECT_ID")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		// TODO: Handle error.
		log.Printf("failed to create pubsub client: %v", err)
	}
	topic := client.Topic(os.Getenv("CANDIDATE_GMAIL_MESSAGES_TOPIC"))
	defer topic.Stop()

	rawMessage, err := json.Marshal(emailMessages)
	if err != nil {
		log.Printf("failed to marshal email messages: %v", err)
	}

	// publish message
	result := topic.Publish(ctx, &pubsub.Message{
		Data: rawMessage,
	})

	// Block until the result is returned and a server-generated ID is returned
	// for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		log.Printf("failed to publish message: %v", err)
	}
	log.Printf("Published a message; msg ID: %v", id)
}

// triggerBackgroundfFullEmailSync triggers a background function to sync all emails since a given date
// Note: The service account must be a principal with invoker permission on the full sync service
func triggerBackgroundfFullEmailSync(ctx context.Context, email string, startDate time.Time) error {
	triggerURL := os.Getenv("TRIGGER_FULL_SYNC_URL")
	if triggerURL == "" {
		return errors.New("TRIGGER_FULL_SYNC_URL is not set")
	}
	httpClient, err := idtoken.NewClient(ctx, triggerURL)
	if err != nil {
		return err
	}

	// trigger full sync
	body := FullEmailSyncRequest{
		Email:     email,
		StartDate: startDate,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// fire and forget
	// trigger in go routine, so that the function can return and the cloudfunction can be scaled down
	// the proper way to do this is to use pubsub for automatic retries, but this is a quick and dirty solution
	go func() {
		log.Printf("trigger full sync for %s from %s", email, startDate.Format("2006/01/02"))
		_, err = httpClient.Post(triggerURL, "application/json", io.NopCloser(bytes.NewReader(bodyBytes)))
		if err != nil {
			log.Printf("failed to trigger full sync: %v", err)
		}
	}()

	return nil
}

func syncNewEmails(
	user db.UserProfile,
	srv *srcmail.Service,
	syncHistory db.UserEmailSyncHistory,
) error {
	var err error
	pageToken := ""
	var messages []*gmail.Message
	historyIDExpired := false

	for {
		// Make request to fetch new emails from previous history id or last sync date
		// get next set of messages
		if historyIDExpired {
			// if history id is expired, trigger async full sync to last sync date
			err = triggerBackgroundfFullEmailSync(context.Background(), user.Email, syncHistory.SyncedAt)
			if err != nil {
				return fmt.Errorf("error triggering full sync: %v", err)
			}
			// done!
			return nil
		} else {
			messages, pageToken, err = fetchNewEmailsSinceHistoryID(srv, uint64(syncHistory.HistoryID), "UNREAD", pageToken)
		}

		// for now, abort on error
		if err != nil {
			// check for a history not found error
			if srcmail.IsNotFound(err) && !historyIDExpired {
				log.Printf("expired history id: %v", err)
				log.Printf("syncing from %s", syncHistory.SyncedAt.Format("2006/01/02"))
				// set flag and continue iterating
				historyIDExpired = true
				continue
			}
			return fmt.Errorf("error fetching emails: %v", err)
		}

		// TODO: Publish all messages before waiting on results
		// var results []*pubsub.PublishResult
		publishMessages(user.Email, messages)

		if pageToken == "" {
			break
		}
	}

	return nil
}
