package cloudfunctions

import (
	"google.golang.org/api/gmail/v1"

	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const maxResults = 250

var (
	// historyTypes is a list of all the history types that we want to fetch.
	// We don't care about deleted messages.
	// See https://developers.google.com/gmail/api/v1/reference/users/history/list
	historyTypes = []string{
		"messageAdded",
		"labelAdded",
		"labelRemoved",
	}
)

func fetchChangesSinceHistoryID(srv *srcmail.Service, historyID uint64, pageToken string) ([]*gmail.History, string, error) {
	r, err := srcmail.ExecuteWithRetries(func() (*gmail.ListHistoryResponse, error) {
		return srv.Users.History.
			List(srv.UserID).
			StartHistoryId(historyID).
			HistoryTypes(historyTypes...).
			PageToken(pageToken).
			MaxResults(maxResults).
			Do()
	})

	if err != nil {
		return nil, "", err
	}

	return r.History, r.NextPageToken, nil
}

func historyToAddedMessages(histories []*gmail.History) []*gmail.Message {
	messages := []*gmail.Message{}
	for _, h := range histories {
		for _, m := range h.MessagesAdded {
			messages = append(messages, m.Message)
		}
	}
	return messages
}

func (cf *CloudFunction) historyToEmailLabelChanges(histories []*gmail.History) *schema.EmailLabelChanges {
	changes := []schema.EmailLabelChange{}
	for _, h := range histories {
		for _, m := range h.LabelsAdded {
			changes = append(changes, schema.EmailLabelChange{
				MessageID:  m.Message.Id,
				LabelIDs:   m.LabelIds,
				ChangeType: schema.EmailLabelChangeTypeAdded,
			})
		}
		for _, m := range h.LabelsRemoved {
			changes = append(changes, schema.EmailLabelChange{
				MessageID:  m.Message.Id,
				LabelIDs:   m.LabelIds,
				ChangeType: schema.EmailLabelChangeTypeRemoved,
			})
		}
	}
	return &schema.EmailLabelChanges{
		Email:   cf.payload.Email,
		Changes: &changes,
	}
}
