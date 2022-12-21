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
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail/label"
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
		log.Printf("trigger full sync for %s from %s", email, startDate.Format("2006/01/02"))
		_, err = httpClient.Post(triggerURL, "application/json", ioutil.NopCloser(bytes.NewReader(bodyBytes)))
		if err != nil {
			log.Printf("failed to trigger full sync: %v", err)
		}
	}()

	return nil
}

func syncNewEmails(
	user client.UserProfile,
	srv *mail.Service,
	queries client.Querier,
	classifier Classifier,
	syncHistory client.UserEmailSyncHistory,
	labels *srclabel.Labels,
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
			gErr := &googleapi.Error{}
			if errors.As(err, &gErr); !historyIDExpired && gErr.Code == http.StatusNotFound {
				log.Printf("expired history id: %v", gErr)
				log.Printf("syncing from %s", syncHistory.SyncedAt.Format("2006/01/02"))
				// set flag and continue iterating
				historyIDExpired = true
				continue
			}
			return fmt.Errorf("error fetching emails: %v", err)
		}

		// process messages
		examples := map[string]*PredictRequest{}
		for _, m := range messages {
			// payload isn't included in the list endpoint responses
			message, err := srv.Users.Messages.Get(srv.UserID, m.Id).Do()
			if err != nil {
				gErr := &googleapi.Error{}
				if errors.As(err, &gErr); gErr.Code == http.StatusNotFound {
					// message was deleted, skip
					log.Printf("skipping message %s was deleted", m.Id)
					continue
				}
				// for now, abort on other errors
				return fmt.Errorf("error getting message %s: %v", m.Id, err)
			}

			// filter out empty messages
			if message.Payload == nil {
				continue
			}

			// filter out messages with the sent or already have a job label
			// TODO: Handle new messages with the job label differently
			if contains(message.LabelIds, "SENT") || contains(message.LabelIds, labels.JobsOpportunity.Id) {
				continue
			}

			// check if message sender is on the allow list
			allowed, err := srv.AllowMessage(message)
			if err != nil {
				log.Printf("error checking allow list: %v", err)
			}
			// do not take action on allowed senders
			if allowed {
				log.Printf("allowing message: %s", message.Id)
				continue
			}

			example := &PredictRequest{
				From:    mail.MessageSender(message),
				Subject: mail.MessageSubject(message),
				Body:    mail.MessageBody(message),
			}
			examples[message.Id] = example
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
		err = handleRecruitingEmails(srv, user, labels, recruitingEmailIDs)
		// for now, abort on error
		if err != nil {
			return fmt.Errorf("error modifying recruiting emails: %v", err)
		}

		// save statistics
		if len(examples) > 0 {
			err = queries.IncrementUserEmailStat(
				context.Background(),
				client.IncrementUserEmailStatParams{
					UserID:    user.UserID,
					Email:     user.Email,
					StatID:    "emails_processed",
					StatValue: int32(len(examples)),
				},
			)
			if err != nil {
				// print error, but don't abort
				log.Printf("error incrementing user email stat: %v", err)
			}
		}
		if len(recruitingEmailIDs) > 0 {
			err = queries.IncrementUserEmailStat(
				context.Background(),
				client.IncrementUserEmailStatParams{
					UserID:    user.UserID,
					Email:     user.Email,
					StatID:    "jobs_detected",
					StatValue: int32(len(recruitingEmailIDs)),
				},
			)
			if err != nil {
				// print error, but don't abort
				log.Printf("error incrementing user email stat: %v", err)
			}
		}

		if pageToken == "" {
			break
		}
	}
	return nil
}

func handleRecruitingEmails(srv *mail.Service, profile client.UserProfile, labels *srclabel.Labels, messageIDs []string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	removeLabels := []string{}
	if profile.AutoArchive {
		removeLabels = append(removeLabels, "INBOX", "UNREAD")
	}

	err := srv.Users.Messages.BatchModify(srv.UserID, &gmail.BatchModifyMessagesRequest{
		Ids: messageIDs,
		// Add job opportunity label and parent folder labels
		AddLabelIds:    []string{labels.SRC.Id, labels.Jobs.Id, labels.JobsOpportunity.Id},
		RemoveLabelIds: removeLabels,
	}).Do()

	if err != nil {
		return fmt.Errorf("error modifying recruiting emails: %v", err)
	}

	if profile.AutoContribute {
		for _, id := range messageIDs {
			// shouldn't happen
			if examplesCollectorSrv == nil {
				log.Print("examples collector service not initialized")
				break
			}
			// clone the message to the examples inbox
			_, err := mail.CloneMessage(srv, examplesCollectorSrv, id, collectedExampleLabels)

			if err != nil {
				// don't abort on error
				log.Printf("error collecting email %s: %v", id, err)
				continue
			}
		}
	}

	return nil
}
