package cloudfunctions

import (
	"fmt"
	"time"

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

// GetEmailsSinceLastYear syncs all emails from the previous year to today.
// It returns the messages fetched, the next page token, and an error.
// Use the next page token to fetch the rest of the emails
func GetEmailsSinceLastYear(srv *gmail.Service, userID, pageToken string) ([]*gmail.Message, string, error) {
	// get the date from one year ago, even archived messages (ignore deleted)
	oneYearAgo := time.Now().AddDate(-1, 0, 0).Format("2006/01/02")
	q := fmt.Sprintf("-label:sent after:%s", oneYearAgo)

	m, err := srv.Users.Messages.List(userID).Q(q).PageToken(pageToken).MaxResults(maxResults).Do()

	return m.Messages, m.NextPageToken, err
}
