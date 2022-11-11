package gmail

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var scopes = []string{
	gmail.GmailModifyScope,
	gmail.GmailLabelsScope,
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
