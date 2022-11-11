package gmail

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/api/gmail/v1"

	"github.com/jaytaylor/html2text"
)

func getHTMLContentFromMessageParts(m *gmail.MessagePart) (string, error) {
	if m.MimeType == "text/html" {
		decoded, err := base64.URLEncoding.DecodeString(m.Body.Data)
		text, err := html2text.FromString(string(decoded),
			html2text.Options{
				PrettyTables: false,
				// text only loses some information about headers, lists, block quotes, bold, italic text
				// https://github.com/jaytaylor/html2text/pull/49/files
				TextOnly: false,
			})

		if err != nil {
			return "", err
		}

		return text, nil
	}

	// recursively look for a html-based body
	for _, p := range m.Parts {
		content, mimeType := getHTMLContentFromMessageParts(p)
		if content != "" {
			return content, mimeType
		}
	}

	return "", nil
}

func getTextContentFromMessageParts(m *gmail.MessagePart) (string, error) {
	if m.MimeType == "text/plain" {
		decoded, err := base64.URLEncoding.DecodeString(m.Body.Data)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}

	// recursively look for a text-based body
	for _, p := range m.Parts {
		content, _ := getTextContentFromMessageParts(p)
		if content != "" {
			return content, nil
		}
	}

	return "", nil
}

func getContentFromMessageParts(m *gmail.MessagePart) (string, error) {
	// try to get native text content first
	// else get html content and convert to text
	text, err := getTextContentFromMessageParts(m)

	// if there is an error or any text, return it
	if (err != nil) || (text != "") {
		return text, err
	}

	return getHTMLContentFromMessageParts(m)
}

func MessageToString(m *gmail.Message) (string, error) {
	var from, subject string

	// get sender and subject from headers
	for _, h := range m.Payload.Headers {
		if h.Name == "Subject" {
			subject = h.Value
		}
		if h.Name == "From" {
			from = h.Value
		}
	}

	// parse body parts for body text
	body, err := getContentFromMessageParts(m.Payload)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(` From: %s
 Subject: %s
 Body: %s`, from, subject, body), nil
}
