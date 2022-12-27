package gmail

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"

	"github.com/cenkalti/backoff/v4"
)

// IsOAuth2Error checks if the error is an oauth2.RetrieveError.
// This is useful for detecting expired or revoked access tokens.
func IsOAuth2Error(err error) bool {
	target := &oauth2.RetrieveError{}
	return errors.As(err, &target)
}

// IsGoogleAPIError checks if the error is a googleapi.Error
func IsGoogleAPIError(err error) bool {
	target := &googleapi.Error{}
	return errors.As(err, &target)
}

// IsRateLimitError checks for a status too many requests (429) response from a Google API
func IsRateLimitError(err error) bool {
	return IsGoogleAPIError(err) && err.(*googleapi.Error).Code == http.StatusTooManyRequests
}

// IsNotFound checks for a status not found (404) response from a Google API
func IsNotFound(err error) bool {
	return IsGoogleAPIError(err) && err.(*googleapi.Error).Code == http.StatusNotFound
}

// ExecuteWithRetries executes a function
// and automatically retries with an exponential back-off if the function rate turns a
// googleapi rate limit error (IsRateLimitError).
func ExecuteWithRetries[T any](f func() (T, error)) (T, error) {
	t, err := f()
	if IsRateLimitError(err) {
		return backoff.RetryWithData(f, backoff.NewExponentialBackOff())
	}
	return t, err
}
