package gmail_test

import (
	"errors"
	"net/http"
	"testing"

	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"

	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
)

func TestIsOAuth2Error(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "oauth2.RetrieveError",
			err:  &oauth2.RetrieveError{},
			want: true,
		},
		{
			name: "non oauth error",
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := srcmail.IsOAuth2Error(tc.err)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestIsGoogleAPIError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "googleapi.Error",
			err:  &googleapi.Error{},
			want: true,
		},
		{
			name: "non googleapi error",
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := srcmail.IsGoogleAPIError(tc.err)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "rate limit error",
			err: &googleapi.Error{
				Code: http.StatusTooManyRequests,
			},
			want: true,
		},
		{
			name: "other googleapi error",
			err: &googleapi.Error{
				Code: http.StatusNotFound,
			},
			want: false,
		},
		{
			name: "non googleapi",
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := srcmail.IsRateLimitError(tc.err)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "not found error",
			err: &googleapi.Error{
				Code: http.StatusNotFound,
			},
			want: true,
		},
		{
			name: "other googleapi error",
			err: &googleapi.Error{
				Code: http.StatusTooManyRequests,
			},
			want: false,
		},
		{
			name: "non googleapi",
			err:  errors.New("test"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := srcmail.IsNotFound(tc.err)
			if got != tc.want {
				t.Fail()
			}
		})
	}
}

func TestExecuteWithRetries(t *testing.T) {
	tests := []struct {
		name string
		// use a closure to keep track of call counts
		f       func() func() (int, error)
		want    int
		wantErr error
	}{
		{
			name: "success",
			f: func() func() (int, error) {
				n := 0
				return func() (int, error) {
					n++
					return n, nil
				}
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "non rate limit error",
			f: func() func() (int, error) {
				n := 0
				return func() (int, error) {
					n++
					return n, &googleapi.Error{
						Code: http.StatusNotFound,
					}
				}
			},
			want: 1,
			wantErr: &googleapi.Error{
				Code: http.StatusNotFound,
			},
		},
		{
			name: "rate limit error",
			f: func() func() (int, error) {
				n := 0
				return func() (int, error) {
					n++
					t.Logf("retrying %d", n)
					if n < 3 {
						return n, &googleapi.Error{
							Code: http.StatusTooManyRequests,
						}
					}
					return n, nil
				}
			},
			want:    3,
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := tc.f()
			got, err := srcmail.ExecuteWithRetries(f)
			if got != tc.want {
				t.Errorf("want %d calls. got %d calls", tc.want, got)
			}
			if (err == nil && tc.wantErr != nil) ||
				(err != nil && tc.wantErr == nil) ||
				(err != nil && err.Error() != tc.wantErr.Error()) {
				t.Errorf("want %v error. got %v error", tc.wantErr, err)
			}
		})
	}
}
