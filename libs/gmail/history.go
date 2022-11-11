package gmail

import (
	"google.golang.org/api/gmail/v1"
)

const maxResults = 500
const historyTypeMessageAdded = "messageAdded"
const defaultLabelID = "INBOX"

func GetNewEmailsSince(srv *gmail.Service, userID string, historyID int64, labelID string) ([]*gmail.Message, error) {
	if labelID == "" {
		labelID = defaultLabelID
	}

	// Fer now, fetch up to 500 emails instead of paginating
	r, err := srv.Users.History.
		List(userID).
		LabelId(labelID).
		StartHistoryId(historyID).
		HistoryTypes(historyTypeMessageAdded).
		MaxResults(maxResults).
		Do()

	if err != nil {
		return nil, err
	}

	// TODO: Paginate!
	// srv.Users.History.List(userID).PageToken(r.NextPageToken)

	emails := []*gmail.Message{}

	for _, h := range r.History {
		// only look at messages added
		for _, m := range h.MessagesAdded {
			messageID := m.Message.Id
			// payload isn't included in the history response
			message, err := srv.Users.Messages.Get(userID, messageID).Do()

			if err != nil {
				// consider skipping
				return nil, err
			}

			// skip empty messages
			if message.Payload == nil {
				continue
			}

			emails = append(emails, message)
		}
	}

	return emails, nil
}
