package gmail

import (
	"encoding/base64"
	"strings"

	"google.golang.org/api/gmail/v1"
)

func MessageHeader(m *gmail.Message, header string) string {
	header = strings.ToLower(header)
	for _, h := range m.Payload.Headers {
		if strings.ToLower(h.Name) == header {
			return h.Value
		}
	}
	return ""
}

func MessageSender(m *gmail.Message) string {
	return MessageHeader(m, "From")
}

func MessageSubject(m *gmail.Message) string {
	return MessageHeader(m, "Subject")
}

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
