package cloudfunctions

import (
	"fmt"
	"time"

	"google.golang.org/api/gmail/v1"
)

const maxResults = 250

func GetSRCEmails(srv *gmail.Service, userID string, labelId string, startDate time.Time, pageToken string) ([]*gmail.Message, string, error) {
	builder := srv.Users.Messages.List(userID).LabelIds(labelId).PageToken(pageToken).MaxResults(maxResults)

	// always ignore sent emails
	q := "-label:sent"
	if !startDate.IsZero() {
		// start the search from the start date
		q = fmt.Sprintf("%s after:%s", q, startDate.Format("2006/01/02"))
	}

	m, err := builder.Q(q).Do()

	return m.Messages, m.NextPageToken, err
}
