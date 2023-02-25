package cloudfunctions

import (
	"testing"
)

func TestNormaizeBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "empty",
			body: "",
			want: "",
		},
		{
			name: "no html",
			body: "hello",
			want: "hello",
		},
		{
			name: "html",
			body: "hello <b>world</b>",
			want: "hello world",
		},
		{
			name: "html with newlines",
			body: "hello <b>world</b>\n\n",
			want: "hello world",
		},
		{
			name: "html with newlines and spaces",
			body: "hello <b>world</b>\n\n ",
			want: "hello world",
		},
		{
			name: "html with newlines and spaces and tabs",
			body: "hello <b>world</b>\n\n \t",
			want: "hello world",
		},
		{
			name: "html with newlines and spaces and tabs and carriage returns",
			body: "hello <b>world</b>\n\n \t\r",
			want: "hello world",
		},
		{
			name: "html with newlines and spaces and tabs and carriage returns and multiple spaces",
			body: "hello <b>world</b>\n\n \t\r  ",
			want: "hello world",
		},
		{
			name: "html with links",
			body: "hello <b>world</b>\n\n \t\r  <a href=\"https://www.google.com\">google</a>",
			want: "hello world google",
		},
		{
			name: "html with links and newlines",
			body: "hello <b>world</b>\n\n \t\r  <a href=\"https://www.google.com\">google</a>\n\n",
			want: "hello world google",
		},
		{
			name: "text with links",
			body: "hello world https://www.google.com",
			want: "hello world",
		},
		{
			name: "text with links and newlines",
			body: "hello world https://www.google.com\n\n",
			want: "hello world",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeBody(tc.body)
			if got != tc.want {
				t.Errorf("normalizeBody(%q) = %q, want %q", tc.body, got, tc.want)
			}
		})
	}
}
