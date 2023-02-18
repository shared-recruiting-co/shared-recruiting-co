package cloudfunctions

import (
	"google.golang.org/api/gmail/v1"

	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
)

const maxResults = 250
const historyTypeMessageAdded = "messageAdded"

func fetchSentEmailsSinceHistoryID(srv *srcmail.Service, historyID uint64, pageToken string) ([]*gmail.Message, string, error) {
	r, err := srcmail.ExecuteWithRetries(func() (*gmail.ListHistoryResponse, error) {
		return srv.Users.History.
			List(srv.UserID).
			LabelId("SENT").
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
