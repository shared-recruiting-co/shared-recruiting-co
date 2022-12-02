package cloudfunctions

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

func syncNewEmails(
	gmailSrv *gmail.Service,
	gmailUser string,
	classifier Classifier,
	historyID uint64,
	jobLabelID string,
) error {
	var messages []*gmail.Message
	var err error
	pageToken := ""

	for {
		// 7. Make Request to Fetch New Emails from Previous History ID
		// get next set of messages
		// if this is the first notification, trigger a one-time sync
		messages, pageToken, err = GetNewEmailsSince(gmailSrv, gmailUser, historyID, "UNREAD", pageToken)

		// for now, abort on error
		if err != nil {
			// check if it's a googleapi error
			gErr := &googleapi.Error{}
			if errors.As(err, &gErr) {
				// if it's an oauth error, mark the user's token as invalid
				if gErr.Code == http.StatusNotFound {
					// TODO: Handle
					// We want to query since the last successful sync (history.UpdatedAt)
					return fmt.Errorf("error fetching emails: expired history id: %v", gErr)
				}
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

		// 8. Batch predict on new emails
		results, err := classifier.PredictBatch(examples)
		if err != nil {
			return fmt.Errorf("error predicting on examples: %v", err)
		}

		// 9. Get IDs of new recruiting emails
		recruitingEmailIDs := []string{}
		for id, result := range results {
			if !result {
				continue
			}
			recruitingEmailIDs = append(recruitingEmailIDs, id)
		}

		log.Printf("number of recruiting emails: %d", len(recruitingEmailIDs))

		// 10. Take action on recruiting emails
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
