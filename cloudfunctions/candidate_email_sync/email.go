package cloudfunctions

import (
	"fmt"
	"time"

	"google.golang.org/api/gmail/v1"

	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
)

const maxResults = 250

// fetchThreadsSinceDate fetches all threads since the start date
// It ignores threads of only sent emails and threads already processed by SRC
func fetchThreadsSinceDate(srv *srcmail.Service, start time.Time, end time.Time, pageToken string) ([]*gmail.Thread, string, error) {
	// get all (including archived) emails after the start date, ignore sent emails and emails already processed by SRC
	q := fmt.Sprintf("-label:sent -label:%s after:%s", srclabel.SRC.Name, start.Format("2006/01/02"))
	if !end.IsZero() {
		q = fmt.Sprintf("%s before:%s", q, end.Format("2006/01/02"))
	}

	r, err := srcmail.ExecuteWithRetries(func() (*gmail.ListThreadsResponse, error) {
		return srv.Users.Threads.
			List(srv.UserID).
			PageToken(pageToken).
			Q(q).
			MaxResults(maxResults).
			Do()
	})

	if err != nil {
		return nil, "", err
	}

	return r.Threads, r.NextPageToken, nil
}

// SkipThread if the messages already labeled with SRC label
func skipThread(messages []*gmail.Message, label string) bool {
	if len(messages) == 0 {
		return true
	}

	// for each message in the thread, check if it has the @src label
	for _, m := range messages {
		if srcmessage.HasLabel(m, label) {
			return true
		}
	}

	return false
}

func filterMessagesAfterReply(messages []*gmail.Message) []*gmail.Message {
	filtered := []*gmail.Message{}
	// ensure messages are sorted by ascending date
	srcmessage.SortByDate(messages)

	for _, m := range messages {
		if srcmessage.IsSent(m) {
			break
		}
		filtered = append(filtered, m)
	}
	return filtered
}
