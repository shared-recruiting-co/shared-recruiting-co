package gmail

import (
	"google.golang.org/api/gmail/v1"
)

// CloneMessage is a convenience function to cloning a message from one inbox to another.
// On success, it will return the new destination message.
// Cloning differs from forwarding because it preserves the original message sender, recipient, and all other headers.
// TODO: Support recipient anonymization (i.e. randomize 'To' header, strings.Replace first and last name).
func CloneMessage(src *Service, dst *Service, messageID string, dstLabelIds []string) (*gmail.Message, error) {
	// Get the raw message from the source account
	msg, err := ExecuteWithRetries(func() (*gmail.Message, error) {
		return src.Users.Messages.Get(src.UserID, messageID).Format("raw").Do()
	})

	if err != nil {
		return nil, err
	}
	// Clear the source specific message properties
	msg.Id = ""
	msg.ThreadId = ""
	// Set the destination labels
	msg.LabelIds = dstLabelIds
	// Insert the message into the destination account
	return ExecuteWithRetries(func() (*gmail.Message, error) {
		return dst.Users.Messages.Insert(dst.UserID, msg).Do()
	})
}
