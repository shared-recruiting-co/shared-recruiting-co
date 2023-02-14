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
//
// https://developers.google.com/gmail/api/guides/handle-errors#resolve_a_403_error_user_rate_limit_exceeded
func IsGoogleAPIError(err error) bool {
	target := &googleapi.Error{}
	return errors.As(err, &target)
}

// IsRateLimitError checks for a status too many requests (429) response from a Google API
func IsRateLimitError(err error) bool {
	if !IsGoogleAPIError(err) {
		return false
	}
	// cast to googleapi.Error
	gErr := err.(*googleapi.Error)
	if gErr.Code == http.StatusTooManyRequests {
		return true
	}
	if gErr.Code == http.StatusForbidden {
		// check Errors for domain usageLimits
		// Note: ErrorItem doesn't contain a Domain field
		for _, e := range gErr.Errors {
			if e.Reason == "rateLimitExceeded" ||
				e.Reason == "dailyLimitExceeded" ||
				e.Reason == "userRateLimitExceeded" {
				return true
			}
		}

	}
	return false
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
