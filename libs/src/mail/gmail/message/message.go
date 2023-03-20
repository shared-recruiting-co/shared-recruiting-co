package message

import (
	"encoding/base64"
	"regexp"
	"sort"
	"strings"
	"time"

	"google.golang.org/api/gmail/v1"
)

var (
	emailRegex = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
)

// SortByDate sorts messages by date received by gmail (ascending)
// The messages are sorted in place.
func SortByDate(messages []*gmail.Message) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].InternalDate < messages[j].InternalDate
	})
}

// Header returns the value of the header with the given name
func Header(m *gmail.Message, header string) string {
	header = strings.ToLower(header)
	for _, h := range m.Payload.Headers {
		if strings.ToLower(h.Name) == header {
			return h.Value
		}
	}
	return ""
}

// Sender returns the sender of the message
// The sender is the email address of the first "From" header
func Sender(m *gmail.Message) string {
	return Header(m, "From")
}

// Recipient returns the sender of the message
// The sender is the email address of the first "To" header
//
// This doesn't support multiple recipients.
func Recipient(m *gmail.Message) string {
	return Header(m, "To")
}

// RecipientEmail returns the email address of the recipient
// It uses Recipient to get the recipient string, and then extracts the email address from it.
//
// Example: "John Doe" <john.doe@example.com>" -> "john.doe@example.com"
func RecipientEmail(m *gmail.Message) string {
	recipient := Recipient(m)
	if recipient == "" {
		return ""
	}
	return emailRegex.FindString(recipient)
}

// SenderEmail returns the email address of the sender
// It uses Sender to get the sender string, and then extracts the email address from it.
//
// Example: "John Doe" <john.doe@example.com>" -> "john.doe@example.com"
func SenderEmail(m *gmail.Message) string {
	sender := Sender(m)
	if sender == "" {
		return ""
	}
	return emailRegex.FindString(sender)
}

// Subject returns the subject of the message
func Subject(m *gmail.Message) string {
	return Header(m, "Subject")
}

// HostMessageID returns the message id of the message as it was received by the email provider.
//
// This a convenience function that calls Header(m, "Message-ID")
func HostMessageID(m *gmail.Message) string {
	return Header(m, "Message-ID")
}

// Body returns the body of the message as a string
//
// It first checks for a text/plain body.
// If none is found, it checks for a text/html body.
func Body(m *gmail.Message) string {
	// try to get native text content first
	body := getTextContentFromMessageParts(m.Payload)

	if body != "" {
		return body
	}

	// else try get html content
	return getHTMLContentFromMessageParts(m.Payload)
}

func getHTMLContentFromMessageParts(m *gmail.MessagePart) string {
	if m.MimeType == "text/html" {
		decoded, err := base64.URLEncoding.DecodeString(m.Body.Data)
		if err != nil {
			return ""
		}

		return string(decoded)
	}

	// recursively look for a html-based body
	for _, p := range m.Parts {
		content := getHTMLContentFromMessageParts(p)
		if content != "" {
			return content
		}
	}

	return ""
}

func getTextContentFromMessageParts(m *gmail.MessagePart) string {
	if m.MimeType == "text/plain" {
		decoded, err := base64.URLEncoding.DecodeString(m.Body.Data)
		if err != nil {
			return ""
		}
		return string(decoded)
	}

	// recursively look for a text-based body
	for _, p := range m.Parts {
		content := getTextContentFromMessageParts(p)
		if content != "" {
			return content
		}
	}

	return ""
}

// HasLabel returns true if the message contains the given label id
func HasLabel(m *gmail.Message, id string) bool {
	if m.LabelIds == nil {
		return false
	}
	for _, l := range m.LabelIds {
		if l == id {
			return true
		}
	}
	return false
}

// IsSent returns true if the message was sent by the current user
//
// There are a number of ways to check if a message was sent by a user.
// This function checks if the message contains the system "SENT" label, which allows us to only fetch the minimal message information (no headers) from a thread.
func IsSent(m *gmail.Message) bool {
	return HasLabel(m, "SENT")
}

// CreatedAt returns the time the message was created in the email provider.
//
// For senders, this is the time the message was sent.
// For recipients, this is the time the message was received.
func CreatedAt(m *gmail.Message) time.Time {
	return time.Unix(m.InternalDate/1000, 0)
}
