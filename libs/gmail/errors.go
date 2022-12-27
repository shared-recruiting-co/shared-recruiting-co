package gmail

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	"google.golang.org/api/googleapi"

	"github.com/cenkalti/backoff/v4"
)

func IsOAuth2Error(err error) bool {
	target := &oauth2.RetrieveError{}
	return errors.As(err, &target)
}

func IsGoogleAPIError(err error) bool {
	target := &googleapi.Error{}
	return errors.As(err, &target)
}

func IsRateLimitError(err error) bool {
	return IsGoogleAPIError(err) && err.(*googleapi.Error).Code == http.StatusTooManyRequests
}

func IsNotFound(err error) bool {
	return IsGoogleAPIError(err) && err.(*googleapi.Error).Code == http.StatusNotFound
}

func ExecuteWithRetries[T any](f func() (T, error)) (T, error) {
	t, err := f()
	if IsRateLimitError(err) {
		return backoff.RetryWithData(f, backoff.NewExponentialBackOff())
	}
	return t, err
}
