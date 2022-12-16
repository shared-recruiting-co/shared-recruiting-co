package gmail

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
)

const (
	// FwdMsgDelimiter is the delimiter used to separate the forwarded message from the original message.
	FwdMsgDelimiter = "---------- Forwarded message ---------"
)

// ForwardMessage represents a message to be forwarded.
type ForwardMessage struct {
	Sender string
	To     string
	Parent *gmail.Message
}

// GetParentHeader returns the value of the header with the given name from the parent message.
// It ignores casing when looking at headers.
func (m ForwardMessage) GetParentHeader(name string) string {
	name = strings.ToLower(name)
	for _, h := range m.Parent.Payload.Headers {
		if strings.ToLower(h.Name) == name {
			return h.Value
		}
	}
	return ""
}

// InReplyTo returns the contents of the "In-Reply-To:" field of the message to which this one is a reply (the "parent message").
// It follows the spec described in https://tools.ietf.org/html/rfc2822#section-3.6.4
//
// The "In-Reply-To:" field will contain the contents of the "Message-ID:" field of the message to which this one is a reply (the "parent message").
// If there is more than one parent message, then the "In-Reply-To:" field will contain the contents of all of the parents' "Message-ID:" fields.
// If there is no "Message-ID:" field in any of the parent messages, then the new message will have no "In-Reply-To:" field.
func (m ForwardMessage) InReplyTo() string {
	return m.GetParentHeader("Message-ID")
}

// References returns the contents of the "References:" field of the message to which this one is a reply (the "parent message").
// It follows the spec described in https://tools.ietf.org/html/rfc2822#section-3.6.4
//
// The "References:" field will contain the contents of the parent's "References:" field
// (if any) followed by the contents of the parent's "Message-ID:" field (if any).
// If the parent message does not contain a "References:" field
// but does have an "In-Reply-To:" field containing a single message identifier, then the "References:" field
// will contain the contents of the parent's "In-Reply-To:" field followed by
// the contents of the parent's "Message-ID:" field (if any).
// If the parent has none of the "References:", "In-Reply-To:", or "Message-ID:" fields,
// then the new message will have no "References:" field.
func (m ForwardMessage) References() string {

	// references will always contain the messageID (if it exists)
	references := m.GetParentHeader("Message-ID")

	if p := m.GetParentHeader("References"); p != "" {
		if references == "" {
			return p
		}
		return fmt.Sprintf("%s %s", p, references)
	}

	// TODO: Check if inReplyTo is a single message identifier
	if p := m.GetParentHeader("inReplyTo"); p != "" {
		if references == "" {
			return p
		}
		return fmt.Sprintf("%s %s", p, references)
	}

	return references
}

// ParentBody returns the body of the parent message.
// Note: We always convert to text/plain. It's simpler to do deal with and is sufficient for our purposes.
func (m ForwardMessage) ParentBody() string {
	// The "Body:" field will contain the contents of the parent's
	// "Body:" field (if any)
	return MessageBody(m.Parent)
}

// Raw returns the raw RFC-822 compliant forwarded message.
func (m ForwardMessage) Raw() string {
	// New email headers
	fwdHeaders := []gmail.MessagePartHeader{
		{
			Name:  "From",
			Value: m.Sender,
		},
		{
			Name:  "To",
			Value: m.To,
		},
		{
			Name:  "Subject",
			Value: m.GetParentHeader("Subject"),
		},
		{
			Name:  "In-Reply-To",
			Value: m.InReplyTo(),
		},
		{
			Name:  "References",
			Value: m.References(),
		},
	}
	// Original email headers
	// Order matters! (for consistency w/ gmail client) "To:" must be last
	parentHeaders := []gmail.MessagePartHeader{
		{
			Name:  "From",
			Value: m.GetParentHeader("From"),
		},
		{
			Name:  "Date",
			Value: m.GetParentHeader("Date"),
		},
		{
			Name:  "Subject",
			Value: m.GetParentHeader("Subject"),
		},
		{
			Name:  "To",
			Value: m.GetParentHeader("To"),
		},
	}
	// 4. Original content
	parentContent := m.ParentBody()

	// 5. Create the raw RFC-822 compliant forwarded message
	raw := ""
	for _, h := range fwdHeaders {
		// only include non-empty headers
		if h.Value != "" {
			raw += fmt.Sprintf("%s: %s\r\n", h.Name, h.Value)
		}
	}
	// add the delimiter
	raw += fmt.Sprintf("\r\n%s\r\n", FwdMsgDelimiter)
	// add the parent headers
	for _, h := range parentHeaders {
		// all headers should be defined
		raw += fmt.Sprintf("%s: %s\r\n", h.Name, h.Value)
	}
	// add the parent content
	raw += fmt.Sprintf("\r\n\r\n%s", parentContent)

	// 6. Encode the message
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// Create a send-able gmail message from a forwarded message
func (m ForwardMessage) Create() *gmail.Message {
	return &gmail.Message{
		Raw:      m.Raw(),
		ThreadId: m.Parent.ThreadId,
	}
}
