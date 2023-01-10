package cloudfunctions

import (
	"fmt"
	"time"

	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail/label"
	"google.golang.org/api/gmail/v1"
)

const maxResults = 250

// fetchEmailsSinceDate syncs all inbound emails from the start date to today.
// It returns the messages fetched, the next page token, and an error.
// Use the next page token to fetch the rest of the emails
func fetchEmailsSinceDate(srv *mail.Service, date time.Time, pageToken string) ([]*gmail.Message, string, error) {
	// get all (including archived) emails after the start date, ignore sent emails and emails already processed by SRC
	q := fmt.Sprintf("-label:sent -label:%s after:%s", srclabel.SRC.Name, date.Format("2006/01/02"))

	r, err := mail.ExecuteWithRetries(func() (*gmail.ListMessagesResponse, error) {
		return srv.Users.Messages.
			List(srv.UserID).
			PageToken(pageToken).
			Q(q).
			MaxResults(maxResults).
			Do()
	})

	if err != nil {
		return nil, "", err
	}

	return r.Messages, r.NextPageToken, nil
}

// fetchThreadsSinceDate fetches all threads since the start date
// It ignores threads of only sent emails and threads already processed by SRC
func fetchThreadsSinceDate(srv *mail.Service, date time.Time, pageToken string) ([]*gmail.Thread, string, error) {
	// get all (including archived) emails after the start date, ignore sent emails and emails already processed by SRC
	q := fmt.Sprintf("-label:sent -label:%s after:%s", srclabel.SRC.Name, date.Format("2006/01/02"))

	r, err := mail.ExecuteWithRetries(func() (*gmail.ListThreadsResponse, error) {
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
func skipThread(messages []*gmail.Message, srcLabelId string) bool {
	if len(messages) == 0 {
		return true
	}

	// for each message in the thread, check if it has the @src label
	for _, m := range messages {
		if mail.MessageHasLabel(m, srcLabelId) {
			return true
		}
	}

	return false
}

func filterMessagesAfterReply(messages []*gmail.Message) []*gmail.Message {
	filtered := []*gmail.Message{}
	// ensure messages are sorted by ascending date
	mail.SortMessagesByDate(messages)

	for _, m := range messages {
		if mail.IsMessageSent(m) {
			break
		}
		filtered = append(filtered, m)
	}
	return filtered
}
