package cloudfunctions

import (
	"encoding/base64"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"

	"google.golang.org/api/gmail/v1"
)

const provider = "google"

func init() {
	functions.HTTP("PopulateJobs", populateJobs)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
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
