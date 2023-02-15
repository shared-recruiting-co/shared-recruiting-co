package cloudfunctions

import (
	"bytes"
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

// triggerBackgroundfFullEmailSync triggers a background function to sync all emails since a given date
// Note: The service account must be a principal with invoker permission on the full sync service
func (cf *CloudFunction) triggerBackgroundfFullEmailSync(startDate time.Time) error {
	triggerURL := os.Getenv("TRIGGER_FULL_SYNC_URL")
	if triggerURL == "" {
		return errors.New("TRIGGER_FULL_SYNC_URL is not set")
	}
	httpClient, err := idtoken.NewClient(cf.ctx, triggerURL)
	if err != nil {
		return err
	}

	// trigger full sync
	body := FullEmailSyncRequest{
		Email:     cf.payload.Email,
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
		log.Printf("trigger full sync for %s from %s", cf.payload.Email, startDate.Format("2006/01/02"))
		_, err = httpClient.Post(triggerURL, "application/json", io.NopCloser(bytes.NewReader(bodyBytes)))
		if err != nil {
			log.Printf("failed to trigger full sync: %v", err)
		}
	}()

	return nil
}

func (cf *CloudFunction) syncHistory(
	syncHistory db.UserEmailSyncHistory,
) error {
	var err error
	pageToken := ""
	var messages []*gmail.Message
	var results []*pubsub.PublishResult
	historyIDExpired := false

	defer cf.topics.CandidateGmailMessages.Stop()

	for {
		// Make request to fetch new emails from previous history id or last sync date
		// get next set of messages
		if historyIDExpired {
			// if history id is expired, trigger async full sync to last sync date
			err = cf.triggerBackgroundfFullEmailSync(syncHistory.SyncedAt)
			if err != nil {
				return fmt.Errorf("error triggering full sync: %v", err)
			}
			// done!
			return nil
		} else {
			messages, pageToken, err = fetchNewEmailsSinceHistoryID(cf.srv, uint64(syncHistory.HistoryID), pageToken)
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

		if len(messages) > 0 {
			result, err := cf.PublishCandidateMessages(messages)
			if err != nil {
				return fmt.Errorf("error publishing candidate messages: %w", err)
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

func (cf *CloudFunction) PublishCandidateMessages(messages []*gmail.Message) (*pubsub.PublishResult, error) {
	emailMessages := schema.EmailMessages{
		Email:    cf.payload.Email,
		Messages: make([]string, len(messages)),
	}
	for i, message := range messages {
		emailMessages.Messages[i] = message.Id
	}
	rawMessage, err := json.Marshal(emailMessages)
	if err != nil {
		return nil, fmt.Errorf("error marshalling email messages: %w", err)
	}

	// publish message
	result := cf.topics.CandidateGmailMessages.Publish(cf.ctx, &pubsub.Message{
		Data: rawMessage,
	})

	return result, nil
}
