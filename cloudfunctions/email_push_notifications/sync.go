package cloudfunctions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/idtoken"
)

type FullEmailSyncRequest struct {
	Email     string    `json:"email"`
	StartDate time.Time `json:"start_date"`
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
		_, err = httpClient.Post(triggerURL, "application/json", ioutil.NopCloser(bytes.NewReader(bodyBytes)))
		if err != nil {
			log.Printf("failed to trigger full sync: %v", err)
		}
	}()

	return nil
}

func syncNewEmails(
	gmailSrv *gmail.Service,
	gmailUser string,
	classifier Classifier,
	syncHistory client.UserEmailSyncHistory,
	jobLabelID string,
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
			err = triggerBackgroundfFullEmailSync(context.Background(), "TODO", syncHistory.SyncedAt.Time)
			if err != nil {
				return fmt.Errorf("error triggering full sync: %v", err)
			}
			// done!
			return nil
		} else {
			messages, pageToken, err = getNewEmailsSinceHistoryID(gmailSrv, gmailUser, uint64(syncHistory.HistoryID), "UNREAD", pageToken)
		}

		// for now, abort on error
		if err != nil {
			// check for a history not found error
			gErr := &googleapi.Error{}
			if errors.As(err, &gErr); !historyIDExpired && gErr.Code == http.StatusNotFound {
				log.Printf("expired history id: %v", gErr)
				// make sure the sync at date is set
				if !syncHistory.SyncedAt.Valid {
					return fmt.Errorf("history id expired, but user has never synced before")
				}
				log.Printf("syncing from %s", syncHistory.SyncedAt.Time.Format("2006/01/02"))
				// set flag and continue iterating
				historyIDExpired = true
				continue
			}
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

			// filter out empty messages
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

		// Batch predict on new emails
		results, err := classifier.PredictBatch(examples)
		if err != nil {
			return fmt.Errorf("error predicting on examples: %v", err)
		}

		// Get IDs of new recruiting emails
		recruitingEmailIDs := []string{}
		for id, result := range results {
			if !result {
				continue
			}
			recruitingEmailIDs = append(recruitingEmailIDs, id)
		}

		log.Printf("number of recruiting emails: %d", len(recruitingEmailIDs))

		// Take action on recruiting emails
		if len(recruitingEmailIDs) > 0 {
			err = gmailSrv.Users.Messages.BatchModify(gmailUser, &gmail.BatchModifyMessagesRequest{
				Ids: recruitingEmailIDs,
				// Add SRC Label
				AddLabelIds: []string{jobLabelID},
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
	return nil
}
