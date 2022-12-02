package cloudfunctions

import (
	"fmt"
	"time"

	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	"google.golang.org/api/gmail/v1"
)

const maxResults = 250
const historyTypeMessageAdded = "messageAdded"
const defaultLabelID = "UNREAD"

func getNewEmailsSinceHistoryID(srv *gmail.Service, userID string, historyID uint64, labelID string, pageToken string) ([]*gmail.Message, string, error) {
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

func getNewEmailsSinceDate(srv *gmail.Service, userID string, date time.Time, labelID string, pageToken string) ([]*gmail.Message, string, error) {
	if labelID == "" {
		labelID = defaultLabelID
	}

	q := fmt.Sprintf("-label:%s after:%s", mail.SRC_JobOpportunityLabel, date.Format("2006/01/02"))

	r, err := srv.Users.Messages.
		List(userID).
		LabelIds(labelID).
		PageToken(pageToken).
		Q(q).
		MaxResults(maxResults).
		Do()

	if err != nil {
		return nil, "", err
	}

	return r.Messages, r.NextPageToken, nil
}
