package main

import (
	"testing"
)

func TestShortenAcountId(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "noop",
			input: "sa-cf-normal-length",
			want:  "sa-cf-normal-length",
		},
		{
			name:  "shorten",
			input: "sa-cf-very-very-long-service-account-name",
			want:  "sa-cf-very-very-long-service-a",
		},
		{
			name:  "replace common names",
			input: "sa-cf-candidate-recruiter-gmail",
			want:  "sa-cf-ca-re-gm",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := shortenAccountId(tc.input)
			if got != tc.want {
				t.Errorf("ShortenAccountId(%q) = %q, want %q", tc.input, got, tc.want)
			}

		})
	}
}
