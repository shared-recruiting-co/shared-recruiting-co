package cloudfunctions

import (
	"google.golang.org/api/gmail/v1"
)

const maxResults = 250
const historyTypeMessageAdded = "messageAdded"
const defaultLabelID = "UNREAD"

func GetNewEmailsSince(srv *gmail.Service, userID string, historyID uint64, labelID string, pageToken string) ([]*gmail.Message, string, error) {
	if labelID == "" {
		labelID = defaultLabelID
	}

	r, err := srv.Users.History.
		List(userID).
		LabelId(labelID).
		StartHistoryId(historyID).
		HistoryTypes(historyTypeMessageAdded).
		PageToken(pageToken).
		MaxResults(maxResults).
		Do()

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
