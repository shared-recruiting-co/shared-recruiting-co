package cloudfunctions

import (
	"fmt"
	"time"

	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	"google.golang.org/api/gmail/v1"
)

const maxResults = 250

func fetchSRCEmails(srv *mail.Service, startDate time.Time, pageToken string) ([]*gmail.Message, string, error) {
	builder := srv.Users.Messages.List(srv.UserID).PageToken(pageToken).MaxResults(maxResults)

	// always ignore sent emails
	q := fmt.Sprintf("-label:sent label:%s", mail.SRCJobOpportunityLabel)
	if !startDate.IsZero() {
		// start the search from the start date
		q = fmt.Sprintf("%s after:%s", q, startDate.Format("2006/01/02"))
	}

	m, err := builder.Q(q).Do()

	return m.Messages, m.NextPageToken, err
}

func fetchNonSRCEmails(srv *mail.Service, startDate time.Time, pageToken string) ([]*gmail.Message, string, error) {
	builder := srv.Users.Messages.List(srv.UserID).PageToken(pageToken).MaxResults(maxResults)

	// always ignore sent emails
	q := fmt.Sprintf("-label:sent -label:%s", mail.SRCJobOpportunityLabel)
	if !startDate.IsZero() {
		// start the search from the start date
		q = fmt.Sprintf("%s after:%s", q, startDate.Format("2006/01/02"))
	}

	m, err := builder.Q(q).Do()

	return m.Messages, m.NextPageToken, err
}
