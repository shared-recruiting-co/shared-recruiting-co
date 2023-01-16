package cloudfunctions

import (
	"google.golang.org/api/gmail/v1"

	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
)

const maxResults = 250
const historyTypeMessageAdded = "messageAdded"
const defaultLabelID = "UNREAD"

func fetchNewEmailsSinceHistoryID(srv *srcmail.Service, historyID uint64, labelID string, pageToken string) ([]*gmail.Message, string, error) {
	if labelID == "" {
		labelID = defaultLabelID
	}

	r, err := srcmail.ExecuteWithRetries(func() (*gmail.ListHistoryResponse, error) {
		return srv.Users.History.
			List(srv.UserID).
			LabelId(labelID).
			StartHistoryId(historyID).
			HistoryTypes(historyTypeMessageAdded).
			PageToken(pageToken).
			MaxResults(maxResults).
			Do()
	})

	if err != nil {
		return nil, "", err
	}

	messages := []*gmail.Message{}
	for _, h := range r.History {
		// only look at messages added
		for _, m := range h.MessagesAdded {
			messages = append(messages, m.Message)
		}
	}

	return messages, r.NextPageToken, nil
}

// skipThread if the thread
// - thread was started by the user or the user has replied
// - thread has a job label
func skipThread(messages []*gmail.Message, labelID string) bool {
	for _, m := range messages {
		if m.LabelIds != nil {
			if srcmessage.HasLabel(m, "SENT") || srcmessage.HasLabel(m, labelID) {
				return true
			}
		}
	}
	return false
}
