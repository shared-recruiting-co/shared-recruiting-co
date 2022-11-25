package cloudfunctions

import (
	"fmt"
	"time"

	"google.golang.org/api/gmail/v1"
)

const maxResults = 250

func GetSRCEmails(srv *gmail.Service, userID string, startDate time.Time, pageToken string) ([]*gmail.Message, string, error) {
	builder := srv.Users.Messages.List(userID).PageToken(pageToken).MaxResults(maxResults)

	if !startDate.IsZero() {
		// start the search from the start date
		q := fmt.Sprintf("after:%s", startDate.Format("2006/01/02"))
		builder.Q(q)
	}

	m, err := builder.Do()

	return m.Messages, m.NextPageToken, err
}
