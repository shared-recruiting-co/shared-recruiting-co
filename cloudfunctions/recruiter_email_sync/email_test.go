package cloudfunctions

import (
	"testing"

	"google.golang.org/api/gmail/v1"
)

func TestSkipThread(t *testing.T) {
	labelID := "Label_1"
	tests := []struct {
		name     string
		messages []*gmail.Message
		want     bool
	}{
		{
			name: "do not skip",
			messages: []*gmail.Message{
				{
					LabelIds:     []string{"INBOX"},
					InternalDate: 1,
				},
				{
					LabelIds:     []string{},
					InternalDate: 2,
				},
				{
					LabelIds:     []string{"INBOX"},
					InternalDate: 3,
				},
			},
			want: false,
		},
		{
			name: "sent from user",
			messages: []*gmail.Message{
				{
					LabelIds:     []string{"UNREAD", "INBOX"},
					InternalDate: 2,
				},
				{
					LabelIds:     []string{"SENT"},
					InternalDate: 1,
				},
				{
					LabelIds:     []string{"UNREAD", "INBOX"},
					InternalDate: 3,
				},
			},
			want: false,
		},
		{
			name: "already synced",
			messages: []*gmail.Message{
				{
					LabelIds:     []string{labelID, "INBOX"},
					InternalDate: 2,
				},
				{
					LabelIds:     []string{"INBOX"},
					InternalDate: 1,
				},
				{
					LabelIds:     []string{"UNREAD", "INBOX"},
					InternalDate: 3,
				},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := skipThread(tc.messages, labelID)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}
