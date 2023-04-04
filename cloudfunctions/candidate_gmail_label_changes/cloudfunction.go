package cloudfunctions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/getsentry/sentry-go"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const (
	provider = "google"
)

func init() {
	functions.CloudEvent("Handler", handler)
}

// naive error handler for now
func handleError(msg string, err error) error {
	err = fmt.Errorf("%s: %w", msg, err)
	sentry.CaptureException(err)
	return err
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

func contains[T comparable](list []T, item T) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}

type CloudFunction struct {
	ctx                  context.Context
	queries              db.Querier
	srv                  *srcmail.Service
	labels               *srclabel.CandidateLabels
	model                ml.Service
	user                 db.UserProfile
	examplesCollectorSrv *srcmail.Service
	payload              *schema.EmailLabelChanges
	addedLabelFuncs      map[string]func(msg *gmail.Message) error
	removedLabelFuncs    map[string]func(msg *gmail.Message) error
}

func NewCloudFunction(ctx context.Context, payload schema.EmailLabelChanges) (*CloudFunction, error) {
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("error parsing GOOGLE_OAUTH2_CREDENTIALScredentials: %w", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetUserProfileByEmail(ctx, payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile by email: %w", err)
	}

	var examplesCollectorSrv *srcmail.Service
	if user.AutoContribute {
		auth, err := jsonFromEnv("EXAMPLES_GMAIL_OAUTH_TOKEN")
		if err != nil {
			return nil, fmt.Errorf("error parsing EXAMPLES_GMAIL_OAUTH_TOKEN: %w", err)
		}
		examplesCollectorSrv, err = srcmail.NewService(ctx, creds, auth)
		if err != nil {
			return nil, fmt.Errorf("error creating examples collector service: %w", err)
		}
	}

	// 2. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Email:    payload.Email,
		Provider: provider,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user oauth token: %w", err)
	}

	// stop early if user token is marked invalid
	if !userToken.IsValid {
		return nil, fmt.Errorf("user token is not valid: %s", payload.Email)
	}

	// 3. Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := srcmail.NewService(ctx, creds, auth)
	if err != nil {
		return nil, fmt.Errorf("error creating gmail service: %w", err)
	}

	// 4. Get or Create SRC Labels
	labels, err := srv.GetOrCreateCandidateLabels()
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
		if srcmail.IsOAuth2Error(err) {
			log.Printf("error oauth error: %v", err)
			// update the user's oauth token
			err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
				UserID:   userToken.UserID,
				Email:    payload.Email,
				Provider: provider,
				Token:    userToken.Token,
				IsValid:  false,
			})
			if err != nil {
				log.Printf("error updating user oauth token: %v", err)
			} else {
				log.Printf("marked user oauth token as invalid")
			}
		}
		return nil, fmt.Errorf("error getting or creating src labels: %w", err)
	}
	// Create recruiting detector client
	mlServiceBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, mlServiceBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating id token source: %w", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error getting id token: %w", err)
	}

	model := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   mlServiceBaseURL,
		AuthToken: idToken.AccessToken,
	})

	// set up label handlers
	// use label IDs as keys
	addedLabelFuncs := map[string]func(msg *gmail.Message) error{
		labels.JobsOpportunity.Id: handleAddedJobOpportunityLabel,
	}
	removedLabelFuncs := map[string]func(msg *gmail.Message) error{
		labels.JobsOpportunity.Id: handleRemovedJobOpportunityLabel,
	}

	return &CloudFunction{
		ctx:                  ctx,
		queries:              queries,
		srv:                  srv,
		labels:               labels,
		model:                model,
		user:                 user,
		examplesCollectorSrv: examplesCollectorSrv,
		payload:              &payload,
		addedLabelFuncs:      addedLabelFuncs,
		removedLabelFuncs:    removedLabelFuncs,
	}, nil
}

func handler(ctx context.Context, e event.Event) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       os.Getenv("FUNCTION_NAME"),
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %v", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	var msg schema.MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return handleError("error parsing Pub/Sub message", err)
	}

	data := msg.Message.Data
	log.Printf("Event: %s", data)

	var payload schema.EmailLabelChanges
	if err := json.Unmarshal(data, &payload); err != nil {
		return handleError("error parsing email messages", err)
	}

	// validate payload
	// for invalid payloads, we don't want to retry
	if payload.Email == "" || len(payload.Changes) == 0 {
		err = fmt.Errorf("received invalid payload: %v", payload)
		log.Print(err)
		sentry.CaptureException(err)
		return nil
	}

	cf, err := NewCloudFunction(ctx, payload)
	if err != nil {
		return handleError("error creating cloud function", err)
	}

	err = cf.processLabelChanges(payload.Changes)
	if err != nil {
		return handleError("error processing messages", err)
	}

	return nil
}

func (cf *CloudFunction) processLabelChanges(changes []schema.EmailLabelChange) error {
	// placeholder logic
	// for now, we'll just log the changes
	for _, change := range changes {
		log.Printf("change: %v", change)
		// 1. get the message
		msg, err := cf.srv.GetMessage(change.MessageID)
		if err != nil {
			log.Printf("error getting message: %v", err)
			continue
		}
		if change.ChangeType == schema.EmailLabelChangeTypeAdded {
			// 2. for each label, check if we have a function for it
			for _, label := range change.LabelIDs {
				// Check the label is still added
				if !contains(msg.LabelIds, label) {
					log.Printf("label %s was removed from message %s before processing", label, change.MessageID)
					continue
				}
				// process
				err = cf.processLabelAdded(msg, label)
				// abort on error so we can retry
				if err != nil {
					return fmt.Errorf("error processing label added: %w", err)
				}
			}
		} else if change.ChangeType == schema.EmailLabelChangeTypeRemoved {
			// Check the label is still removed
			for _, label := range change.LabelIDs {
				if contains(msg.LabelIds, label) {
					log.Printf("label %s was added back to message %s before processing", label, change.MessageID)
					continue
				}
				// process
				err = cf.processLabelRemoved(msg, label)
				// abort on error so we can retry
				if err != nil {
					return fmt.Errorf("error processing label removed: %w", err)
				}
			}
		} else {
			log.Printf("unknown change type: %v", change.ChangeType)
		}
	}

	return nil
}

func (cf *CloudFunction) processLabelAdded(msg *gmail.Message, labelID string) error {
	// 1. check if we have a function for the label
	fn, ok := cf.addedLabelFuncs[labelID]
	if !ok {
		log.Printf("no add function found for label %s", labelID)
		return nil
	}
	// 2. run the function
	return fn(msg)
}

func (cf *CloudFunction) processLabelRemoved(msg *gmail.Message, labelID string) error {
	// 1. check if we have a function for the label
	fn, ok := cf.removedLabelFuncs[labelID]
	if !ok {
		log.Printf("no remove function found for label %s", labelID)
		return nil
	}
	// 2. run the function
	return fn(msg)
}

// Label Functions

func handleAddedJobOpportunityLabel(msg *gmail.Message) error {
	// For now log
	log.Printf("added job opportunities label: %s", msg.Id)
	// TODO
	// process as a recruiting email
	return nil
}

func handleRemovedJobOpportunityLabel(msg *gmail.Message) error {
	// For now log
	log.Printf("removed job opportunities label: %s", msg.Id)
	// TODO
	// Remove from a user's user_email_job if it exists
	return nil
}
