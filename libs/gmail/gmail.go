package gmail

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail/label"
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

func (s *Service) CreateLabel(l *gmail.Label) (*gmail.Label, error) {
	return s.Users.Labels.Create(s.UserID, l).Do()
}

// GetOrCreateSRCLabels fetches or creates all of the labels managed by SRC
// Label IDs are unique to each gmail account, so we need to list all labels and match based off names.
func (s *Service) GetOrCreateSRCLabels() (*srclabel.Labels, error) {
	// TODO: What is a better data structure or interface to make this function more DRY?
	// list all labels
	labels, err := s.Users.Labels.List(s.UserID).Do()
	if err != nil {
		return nil, err
	}
	result := srclabel.Labels{}

	// for each label,
	// if it exists, add it to the struct
	// if it doesn't exist, create it and add it to the struct
	for _, label := range labels.Labels {
		switch label.Name {
		case srclabel.SRC.Name:
			result.SRC = label
			break
		case srclabel.Jobs.Name:
			result.Jobs = label
			break
		case srclabel.JobsOpportunity.Name:
			result.JobsOpportunity = label
			break
		case srclabel.Allow.Name:
			result.Allow = label
			break
		case srclabel.AllowSender.Name:
			result.AllowSender = label
			break
		case srclabel.AllowDomain.Name:
			result.AllowDomain = label
			break
		case srclabel.Block.Name:
			result.Block = label
			break
		case srclabel.BlockSender.Name:
			result.BlockSender = label
			break
		case srclabel.BlockDomain.Name:
			result.BlockDomain = label
			break
		case srclabel.BlockGraveyard.Name:
			result.BlockGraveyard = label
			break
		}
	}

	// create any labels that don't exist
	if result.SRC == nil {
		result.SRC, err = s.CreateLabel(&srclabel.SRC)
		if err != nil {
			return nil, err
		}
	}
	if result.Jobs == nil {
		result.Jobs, err = s.CreateLabel(&srclabel.Jobs)
		if err != nil {
			return nil, err
		}
	}
	if result.JobsOpportunity == nil {
		result.JobsOpportunity, err = s.CreateLabel(&srclabel.JobsOpportunity)
		if err != nil {
			return nil, err
		}
	}
	if result.Allow == nil {
		result.Allow, err = s.CreateLabel(&srclabel.Allow)
		if err != nil {
			return nil, err
		}
	}
	if result.AllowSender == nil {
		result.AllowSender, err = s.CreateLabel(&srclabel.AllowSender)
		if err != nil {
			return nil, err
		}
	}
	if result.AllowDomain == nil {
		result.AllowDomain, err = s.CreateLabel(&srclabel.AllowDomain)
		if err != nil {
			return nil, err
		}
	}
	if result.Block == nil {
		result.Block, err = s.CreateLabel(&srclabel.Block)
		if err != nil {
			return nil, err
		}

	}
	if result.BlockSender == nil {
		result.BlockSender, err = s.CreateLabel(&srclabel.BlockSender)
		if err != nil {
			return nil, err
		}
	}
	if result.BlockDomain == nil {
		result.BlockDomain, err = s.CreateLabel(&srclabel.BlockDomain)
		if err != nil {
			return nil, err
		}
	}
	if result.BlockGraveyard == nil {
		result.BlockGraveyard, err = s.CreateLabel(&srclabel.BlockGraveyard)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
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

// GetMessage() fetches a message by ID
// It is a convenience method that wraps the gmail API call.
// It opens opportunities for caching and rate-limiting handling
func (s *Service) GetMessage(id string) (*gmail.Message, error) {
	return s.Users.Messages.Get(s.UserID, id).Do()
}

// Ideas for future
// ListMessages()
// cache GetMessage calls until next ListMessages call
// Something to help with setting the default query for listing src and non-src messages
