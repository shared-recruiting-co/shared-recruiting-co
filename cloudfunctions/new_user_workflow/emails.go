package cloudfunctions

import (
	"fmt"
	"time"

	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	"google.golang.org/api/gmail/v1"
)

const maxResults = 250

// GetEmailsSinceLastYear syncs all emails from the previous year to today.
// It returns the messages fetched, the next page token, and an error.
// Use the next page token to fetch the rest of the emails
func GetEmailsSinceLastYear(srv *gmail.Service, userID, pageToken string) ([]*gmail.Message, string, error) {
	// get the date from one year ago, even archived messages (ignore deleted)
	oneYearAgo := time.Now().AddDate(-1, 0, 0).Format("2006/01/02")
	q := fmt.Sprintf("-label:sent -label:%s after:%s", mail.SRC_JobOpportunityLabel, oneYearAgo)

	m, err := srv.Users.Messages.List(userID).Q(q).PageToken(pageToken).MaxResults(maxResults).Do()

	return m.Messages, m.NextPageToken, err
}
