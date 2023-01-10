package gmail

import (
	"encoding/base64"
	"sort"
	"strings"

	"google.golang.org/api/gmail/v1"
)

// SortMessagesByDate sorts messages by date received by gmail (ascending)
// The messages are sorted in place.
func SortMessagesByDate(messages []*gmail.Message) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].InternalDate < messages[j].InternalDate
	})
}

// MessageHeader returns the value of the header with the given name
func MessageHeader(m *gmail.Message, header string) string {
	header = strings.ToLower(header)
	for _, h := range m.Payload.Headers {
		if strings.ToLower(h.Name) == header {
			return h.Value
		}
	}
	return ""
}

// MessageSender returns the sender of the message
// The sender is the email address of the first "From" header
func MessageSender(m *gmail.Message) string {
	return MessageHeader(m, "From")
}

// MessageSubject returns the subject of the message
func MessageSubject(m *gmail.Message) string {
	return MessageHeader(m, "Subject")
}

// MessageBody returns the body of the message as a string
//
// It first checks for a text/plain body.
// If none is found, it checks for a text/html body.
func MessageBody(m *gmail.Message) string {
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

// MessageHasLabel returns true if the message contains the given label id
func MessageHasLabel(m *gmail.Message, id string) bool {
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

// IsMessageSent returns true if the message was sent by the current user
//
// There are a number of ways to check if a message was sent by a user.
// This function checks if the message contains the system "SENT" label, which allows us to only fetch the minimal message information (no headers) from a thread.
func IsMessageSent(m *gmail.Message) bool {
	return MessageHasLabel(m, "SENT")
}
