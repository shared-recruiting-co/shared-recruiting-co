package message_test

import (
	"encoding/base64"
	"testing"
	"time"

	message "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"google.golang.org/api/gmail/v1"
)

func TestSortMessagesByDate(t *testing.T) {
	tests := []struct {
		name     string
		messages []*gmail.Message
		want     []*gmail.Message
	}{
		{
			name: "Sort",
			messages: []*gmail.Message{
				{
					Id:           "2",
					InternalDate: 2,
				},
				{
					Id:           "1",
					InternalDate: 1,
				},
			},
			want: []*gmail.Message{
				{
					Id:           "1",
					InternalDate: 1,
				},
				{
					Id:           "2",
					InternalDate: 2,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			message.SortByDate(tc.messages)
			for i, m := range tc.messages {
				if m.Id != tc.want[i].Id {
					t.Errorf("got %s, want %s", m.Id, tc.want[i].Id)
				}
			}
		})
	}
}

func TestMessageHeader(t *testing.T) {
	msg := &gmail.Message{
		Payload: &gmail.MessagePart{
			Headers: []*gmail.MessagePartHeader{
				{
					Name:  "From",
					Value: "from",
				},
				{
					Name:  "TO",
					Value: "to",
				},
				{
					Name:  "Subject",
					Value: "subject",
				},
			},
		},
	}

	tests := []struct {
		name   string
		header string
		want   string
	}{
		{
			name:   "Simple",
			header: "From",
			want:   "from",
		},
		{
			name:   "Case Insensitive",
			header: "TO",
			want:   "to",
		},
		{
			name:   "Missing",
			header: "Missing",
			want:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.Header(msg, tc.header)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestMessageBody(t *testing.T) {
	tests := []struct {
		name    string
		message *gmail.Message
		want    string
	}{
		{
			name: "text/plain",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					MimeType: "text/plain",
					Body: &gmail.MessagePartBody{
						Data: base64.URLEncoding.EncodeToString([]byte("data")),
					},
				},
			},
			want: "data",
		},
		{
			name: "text/html",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					MimeType: "text/html",
					Body: &gmail.MessagePartBody{
						Data: base64.URLEncoding.EncodeToString([]byte("<html>data</html>")),
					},
				},
			},
			want: "<html>data</html>",
		},
		{
			name: "No Body",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{},
			},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.Body(tc.message)
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestMessageSenderEmail(t *testing.T) {
	tests := []struct {
		name    string
		message *gmail.Message
		want    string
	}{
		{
			name: "Only Email",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "from",
							Value: "test@example.com",
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Display Name",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "From",
							Value: "Test Test <test@example.com>",
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Whitespace",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "from",
							Value: "Test Test < test@example.com>",
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Multiple ><",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "FROM",
							Value: `"ANNOUNCEMENT <<>>" <test@example.com>`,
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Empty",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "from",
							Value: "",
						},
					},
				},
			},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.SenderEmail(tc.message)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestMessageRecipientEmail(t *testing.T) {
	tests := []struct {
		name    string
		message *gmail.Message
		want    string
	}{
		{
			name: "Only Email",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "to",
							Value: "test@example.com",
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Display Name",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "to",
							Value: "Test Test <test@example.com>",
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Whitespace",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "To",
							Value: "Test Test < test@example.com>",
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Multiple ><",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "To",
							Value: `"ANNOUNCEMENT <<>>" <test@example.com>`,
						},
					},
				},
			},
			want: "test@example.com",
		},
		{
			name: "Empty",
			message: &gmail.Message{
				Payload: &gmail.MessagePart{
					Headers: []*gmail.MessagePartHeader{
						{
							Name:  "To",
							Value: "",
						},
					},
				},
			},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.RecipientEmail(tc.message)
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestMessageHasLabel(t *testing.T) {
	label := "label"
	msg := &gmail.Message{
		LabelIds: []string{"one", label, "two"},
	}
	tests := []struct {
		name  string
		label string
		want  bool
	}{
		{
			name:  "Has Label",
			label: label,
			want:  true,
		},
		{
			name:  "Missing Label",
			label: "does not exist",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.HasLabel(msg, tc.label)
			if got != tc.want {
				t.Errorf("got %t, want %t", got, tc.want)
			}
		})
	}
}

func TestIsMessageSent(t *testing.T) {
	tests := []struct {
		name    string
		message *gmail.Message
		want    bool
	}{
		{
			name: "Sent",
			message: &gmail.Message{
				LabelIds: []string{"INBOX", "SENT", "UNREAD"},
			},
			want: true,
		},
		{
			name: "Received",
			message: &gmail.Message{
				LabelIds: []string{"INBOX", "UNREAD"},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.IsSent(tc.message)
			if got != tc.want {
				t.Errorf("got %t, want %t", got, tc.want)
			}
		})
	}
}

func TestCreatedAt(t *testing.T) {
	tests := []struct {
		name    string
		message *gmail.Message
		want    time.Time
	}{
		{
			name: "Sent",
			message: &gmail.Message{
				InternalDate: 123456789,
			},
			want: time.Unix(123456789/1000, 0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := message.CreatedAt(tc.message)
			if got.String() != tc.want.String() {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
