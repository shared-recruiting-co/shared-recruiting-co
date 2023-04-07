// Package gmail is a wrapper Go library for working with the Google gmail API.
// It is in active development and makes no guarantees about API stability. The long-term goal is to create a standard mail client interface we can use across all email clients, like Gmail and Outlook.

package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
)

// scopes required for the SRC service
var scopes = []string{
	gmail.GmailModifyScope,
}

var (
	// senderEmailRe is a regex to extract the email address from a sender string (i.e. "Jo Smo <jo.smo@gmail.com>")
	senderEmailRe = regexp.MustCompile(`<([^>]+)>`)
)

func newGmailService(ctx context.Context, creds []byte, auth []byte) (*gmail.Service, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal(auth, tok)
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}
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

	profile, err := ExecuteWithRetries(func() (*gmail.Profile, error) {
		return s.Users.GetProfile(s.UserID).Do()
	})

	if err != nil {
		return nil, err
	}
	s.profile = profile

	return profile, nil
}

// Labels

func (s *Service) ListLabels() (*gmail.ListLabelsResponse, error) {
	return ExecuteWithRetries(func() (*gmail.ListLabelsResponse, error) {
		return s.Users.Labels.List(s.UserID).Do()
	})
}

func (s *Service) GetOrCreateLabel(l *gmail.Label) (*gmail.Label, error) {
	labels, err := s.ListLabels()

	if err != nil {
		return nil, err
	}

	for _, label := range labels.Labels {
		if label.Name == l.Name {
			return label, err
		}
	}

	return s.CreateLabel(l)
}

func (s *Service) CreateLabel(l *gmail.Label) (*gmail.Label, error) {
	return ExecuteWithRetries(func() (*gmail.Label, error) {
		return s.Users.Labels.Create(s.UserID, l).Do()
	})
}

// GetOrCreateCandidateLabels fetches or creates all of the labels managed by SRC
// Label IDs are unique to each gmail account, so we need to list all labels and match based off names.
func (s *Service) GetOrCreateCandidateLabels() (*srclabel.CandidateLabels, error) {
	// TODO: What is a better data structure or interface to make this function more DRY?
	// list all labels
	labels, err := s.ListLabels()
	if err != nil {
		return nil, err
	}
	result := srclabel.CandidateLabels{}

	// for each label,
	// if it exists, add it to the struct
	// if it doesn't exist, create it and add it to the struct
	for _, label := range labels.Labels {
		switch label.Name {
		case srclabel.SRC.Name:
			result.SRC = label
		case srclabel.Jobs.Name:
			result.Jobs = label
		case srclabel.JobsOpportunity.Name:
			result.JobsOpportunity = label
		case srclabel.JobsInterested.Name:
			result.JobsInterested = label
		case srclabel.JobsNotInterested.Name:
			result.JobsNotInterested = label
		case srclabel.JobsSaved.Name:
			result.JobsSaved = label
		case srclabel.Allow.Name:
			result.Allow = label
		case srclabel.AllowSender.Name:
			result.AllowSender = label
		case srclabel.AllowDomain.Name:
			result.AllowDomain = label
		case srclabel.Block.Name:
			result.Block = label
		case srclabel.BlockSender.Name:
			result.BlockSender = label
		case srclabel.BlockDomain.Name:
			result.BlockDomain = label
		case srclabel.BlockGraveyard.Name:
			result.BlockGraveyard = label
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
	if result.JobsInterested == nil {
		result.JobsInterested, err = s.CreateLabel(&srclabel.JobsInterested)
		if err != nil {
			return nil, err
		}
	}
	if result.JobsNotInterested == nil {
		result.JobsNotInterested, err = s.CreateLabel(&srclabel.JobsNotInterested)
		if err != nil {
			return nil, err
		}
	}
	if result.JobsSaved == nil {
		result.JobsSaved, err = s.CreateLabel(&srclabel.JobsSaved)
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

// GetOrCreateRecruiterLabels fetches or creates all of the labels managed by SRC
// Label IDs are unique to each gmail account, so we need to list all labels and match based off names.
func (s *Service) GetOrCreateRecruiterLabels() (*srclabel.RecruiterLabels, error) {
	// TODO: What is a better data structure or interface to make this function more DRY?
	// list all labels
	labels, err := s.ListLabels()
	if err != nil {
		return nil, err
	}
	result := srclabel.RecruiterLabels{}

	// for each label,
	// if it exists, add it to the struct
	// if it doesn't exist, create it and add it to the struct
	for _, label := range labels.Labels {
		switch label.Name {
		case srclabel.SRC.Name:
			result.SRC = label
		case srclabel.Recruiting.Name:
			result.Recruiting = label
		case srclabel.RecruitingOutbound.Name:
			result.RecruitingOutbound = label
		}
	}

	// create any labels that don't exist
	if result.SRC == nil {
		result.SRC, err = s.CreateLabel(&srclabel.SRC)
		if err != nil {
			return nil, err
		}
	}
	if result.Recruiting == nil {
		result.Recruiting, err = s.CreateLabel(&srclabel.Recruiting)
		if err != nil {
			return nil, err
		}
	}
	if result.RecruitingOutbound == nil {
		result.RecruitingOutbound, err = s.CreateLabel(&srclabel.RecruitingOutbound)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

// IsSenderAllowed checks if message is on the user's allow list
func (s *Service) IsSenderAllowed(sender string) (bool, error) {
	// if there is a sender display name ("John Doe <john.doe@gmail.com>")
	// extract the email address
	if senderEmailRe.MatchString(sender) {
		matches := senderEmailRe.FindStringSubmatch(sender)
		if len(matches) < 2 {
			return false, fmt.Errorf("error allowing message: unable to parse email from sender: %s", sender)
		}
		sender = matches[1]
	}
	// check if sender is allowed
	// leverage Gmail's native query engine to check
	q := fmt.Sprintf("from:%s label:%s", sender, srclabel.AllowSender.Name)
	resp, err := ExecuteWithRetries(func() (*gmail.ListMessagesResponse, error) {
		return s.Users.Messages.List(s.UserID).Q(q).Do()
	})
	if err != nil {
		return false, err
	}
	if len(resp.Messages) > 0 {
		return true, nil
	}
	// check if sender's email domain is allowed
	parts := strings.SplitAfter(sender, "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("error allowing message: unable to parse domain from sender: %s", sender)
	}
	// remove name related characters (i.e Jo Smo <joe@smo.com>)
	domain := strings.ReplaceAll("@"+parts[1], ">", "")
	q = fmt.Sprintf("from:%s label:%s", domain, srclabel.AllowDomain.Name)
	resp, err = ExecuteWithRetries(func() (*gmail.ListMessagesResponse, error) {
		return s.Users.Messages.List(s.UserID).Q(q).Do()
	})

	if err != nil {
		return false, err
	}

	return len(resp.Messages) > 0, err
}

// IsSenderBlocked checks if message is on the user's block list
func (s *Service) IsSenderBlocked(sender string) (bool, error) {
	// if there is a sender display name ("John Doe <john.doe@gmail.com>")
	// extract the email address
	if senderEmailRe.MatchString(sender) {
		matches := senderEmailRe.FindStringSubmatch(sender)
		if len(matches) < 2 {
			return false, fmt.Errorf("error blocking message: unable to parse email from sender: %s", sender)
		}
		sender = matches[1]
	}
	// check if sender is blocked
	// leverage Gmail's native query engine to check
	q := fmt.Sprintf("from:%s label:%s", sender, srclabel.BlockSender.Name)
	resp, err := ExecuteWithRetries(func() (*gmail.ListMessagesResponse, error) {
		return s.Users.Messages.List(s.UserID).Q(q).Do()
	})
	if err != nil {
		return false, err
	}
	if len(resp.Messages) > 0 {
		return true, nil
	}
	// check if sender's email domain is blocked
	parts := strings.SplitAfter(sender, "@")
	if len(parts) != 2 {
		return false, fmt.Errorf("error blocking message: unable to parse domain from sender: %s", sender)
	}
	// remove name related characters (i.e Jo Smo <joe@smo.com>)
	domain := strings.ReplaceAll("@"+parts[1], ">", "")
	q = fmt.Sprintf("from:%s label:%s", domain, srclabel.BlockDomain.Name)
	resp, err = ExecuteWithRetries(func() (*gmail.ListMessagesResponse, error) {
		return s.Users.Messages.List(s.UserID).Q(q).Do()
	})
	if err != nil {
		return false, err
	}

	return len(resp.Messages) > 0, err
}

// Messages

// BlockMessage blocks a message by moving moving out of the users inbox and into the block graveyard
func (s *Service) BlockMessage(id string, labels *srclabel.CandidateLabels) error {
	_, err := ExecuteWithRetries(func() (*gmail.Message, error) {
		return s.Users.Messages.Modify(s.UserID, id, &gmail.ModifyMessageRequest{
			AddLabelIds:    []string{labels.SRC.Id, labels.Block.Id, labels.BlockGraveyard.Id},
			RemoveLabelIds: []string{"UNREAD", "INBOX"},
		}).Do()
	})
	return err
}

// ForwardEmail is a good enough implementation of forwarding an email in the same format as the gmail client.
// It is good enough because it doesn't naively handle HTML mime-type content or when there are multiple parent messages.
// This is sufficient for our purposes.
func (s *Service) ForwardEmail(messageID, to string) (*gmail.Message, error) {
	f := func() (*gmail.Message, error) {
		// 1. get the original message
		parent, err := s.GetMessage(messageID)
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
	return ExecuteWithRetries(f)
}

// GetMessage() fetches a message by ID
// It is a convenience method that wraps the gmail API call.
// It opens opportunities for caching and rate-limiting handling
func (s *Service) GetMessage(id string) (*gmail.Message, error) {
	f := func() (*gmail.Message, error) {
		return s.Users.Messages.Get(s.UserID, id).Do()
	}
	return ExecuteWithRetries(f)
}

// Threads

// GetThread() fetches a thread by ID with messages in the given format
// See https://developers.google.com/gmail/api/reference/rest/v1/users.threads/get#Format for format values
//
// It is a convenience method that wraps the gmail API call
// It opens opportunities for caching and rate-limiting handling
func (s *Service) GetThread(id, format string) (*gmail.Thread, error) {
	f := func() (*gmail.Thread, error) {
		return s.Users.Threads.Get(s.UserID, id).Format(format).Do()
	}
	return ExecuteWithRetries(f)
}

// Ideas for future
// ListMessages()
// cache GetMessage calls until next ListMessages call
// Something to help with setting the default query for listing src and non-src messages
