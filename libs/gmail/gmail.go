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

// scopes required for the SRC service
var scopes = []string{
	gmail.GmailModifyScope,
}

func newGmailService(ctx context.Context, creds []byte, auth []byte) (*gmail.Service, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal(auth, tok)

	config, err := google.ConfigFromJSON(creds, scopes...)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, tok)

	return gmail.NewService(ctx, option.WithHTTPClient(client))
}

// newDefaultGmailService creates a new gmail service with default credentials on the local machine.
// The default credentials for service accounts do not have a client secret, so they will not auto-refresh the token.
// The auth token must be valid for the lifetime of the service.
func newDefaultGmailService(ctx context.Context, auth []byte) (*gmail.Service, error) {
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

// Service is a wrapper around the gmail service that adds some convenience methods.
type Service struct {
	*gmail.Service
	UserID  string
	profile *gmail.Profile
}

// Create

func NewService(ctx context.Context, creds []byte, auth []byte) (*Service, error) {
	srv, err := newGmailService(ctx, creds, auth)
	if err != nil {
		return nil, err
	}
	return &Service{Service: srv, UserID: "me"}, nil
}

func NewDefaultService(ctx context.Context, auth []byte) (*Service, error) {
	srv, err := newDefaultGmailService(ctx, auth)
	if err != nil {
		return nil, err
	}
	// only set UserID to 'me' for now until we have a use-case for non-authenticated users
	return &Service{Service: srv, UserID: "me"}, nil
}

func (s *Service) Profile() (*gmail.Profile, error) {
	// cache the profile
	if s.profile != nil {
		return s.profile, nil
	}

	profile, err := s.Users.GetProfile(s.UserID).Do()
	if err != nil {
		return nil, err
	}
	s.profile = profile

	return profile, nil
}

// Labels

func (s *Service) GetOrCreateLabel(name string, color *gmail.LabelColor) (*gmail.Label, error) {
	labels, err := s.Users.Labels.List(s.UserID).Do()

	if err != nil {
		return nil, err
	}

	for _, label := range labels.Labels {
		if label.Name == name {
			return label, err
		}
	}

	return s.Users.Labels.Create(s.UserID, &gmail.Label{Name: name, Color: color}).Do()
}

func (s *Service) GetOrCreateSRCLabel() (*gmail.Label, error) {
	return s.GetOrCreateLabel(SRCLabel, SRCLabelColor)
}

func (s *Service) GetOrCreateSRCJobOpportunityLabel() (*gmail.Label, error) {
	return s.GetOrCreateLabel(SRCJobOpportunityLabel, SRCJobOpportunityLabelColor)
}

// Messages

// ForwardEmail is a good enough implementation of forwarding an email in the same format as the gmail client.
// It is good enough because it doesn't naively handle HTML mime-type content or when there are multiple parent messages.
// This is sufficient for our purposes.
func (s *Service) ForwardEmail(messageID, to string) (*gmail.Message, error) {
	// 1. get the original message
	parent, err := s.Users.Messages.Get(s.UserID, messageID).Do()
	if err != nil {
		return nil, err
	}

	// 2. Get the current user's email address
	profile, err := s.Profile()
	if err != nil {
		return nil, err
	}

	// 3. Create the forwarded message
	fwd := ForwardMessage{
		Sender: profile.EmailAddress,
		To:     to,
		Parent: parent,
	}

	// send the message
	return s.Users.Messages.Send(s.UserID, fwd.Create()).Do()
}

// Ideas for future
// GetMessage()
// Something to help with setting the default query for listing src and non-src messages
