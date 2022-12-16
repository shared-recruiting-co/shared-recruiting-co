package cloudfunctions

import (
	"fmt"
	"time"

	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	"google.golang.org/api/gmail/v1"
)

const maxResults = 250

// fetchEmailsSinceDate syncs all inbound emails from the start date to today.
// It returns the messages fetched, the next page token, and an error.
// Use the next page token to fetch the rest of the emails
func fetchEmailsSinceDate(srv *mail.Service, date time.Time, pageToken string) ([]*gmail.Message, string, error) {
	// get all (including archived) emails after the start date, ignore sent emails and emails already processed by SRC
	q := fmt.Sprintf("-label:sent -label:%s after:%s", mail.SRCJobOpportunityLabel, date.Format("2006/01/02"))

	r, err := srv.Users.Messages.
		List(srv.UserID).
		PageToken(pageToken).
		Q(q).
		MaxResults(maxResults).
		Do()

	if err != nil {
		return nil, "", err
	}

	return r.Messages, r.NextPageToken, nil
}
