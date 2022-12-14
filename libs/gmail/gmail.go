package gmail

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var scopes = []string{
	gmail.GmailModifyScope,
}

func NewGmailService(ctx context.Context, creds []byte, auth []byte) (*gmail.Service, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal(auth, tok)

	config, err := google.ConfigFromJSON(creds, scopes...)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, tok)

	return gmail.NewService(ctx, option.WithHTTPClient(client))
}

// NewDefaultGmailService creates a new gmail service with default credentials on the local machine.
// The default credentials for service accounts do not have a client secret, so they will not auto-refresh the token.
// The auth token must be valid for the lifetime of the service.
func NewDefaultGmailService(ctx context.Context, auth []byte) (*gmail.Service, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal(auth, tok)

	if err != nil {
		return nil, err
	}

	ts, err := google.DefaultTokenSource(ctx, scopes...)
	trans := &oauth2.Transport{
		Source: oauth2.ReuseTokenSource(tok, ts),
		Base:   http.DefaultTransport,
	}
	httpClient := &http.Client{Transport: trans}
	return gmail.NewService(ctx, option.WithHTTPClient(httpClient))
}
