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
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
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

type CloudFunction struct {
	ctx      context.Context
	queries  db.Querier
	srv      *srcmail.Service
	model    ml.Service
	user     db.GetRecruiterByEmailRow
	settings schema.EmailMessagesSettings
}

func NewCloudFunction(ctx context.Context, payload schema.EmailMessages) (*CloudFunction, error) {
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("error parsing GOOGLE_OAUTH2_CREDENTIALScredentials: %w", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetRecruiterByEmail(ctx, payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile by email: %w", err)
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

	// Create recruiting detector client
	classifierBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, classifierBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating id token source: %w", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error getting id token: %w", err)
	}

	model := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   classifierBaseURL,
		AuthToken: idToken.AccessToken,
	})

	return &CloudFunction{
		ctx:      ctx,
		queries:  queries,
		srv:      srv,
		model:    model,
		user:     user,
		settings: payload.Settings,
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

	var payload schema.EmailMessages
	if err := json.Unmarshal(data, &payload); err != nil {
		return handleError("error parsing email messages", err)
	}

	// validate payload
	// for invalid payloads, we don't want to retry
	if payload.Email == "" || len(payload.Messages) == 0 {
		err = fmt.Errorf("received invalid payload: %v", payload)
		log.Print(err)
		sentry.CaptureException(err)
		return nil
	}

	_, err = NewCloudFunction(ctx, payload)
	if err != nil {
		return handleError("error creating cloud function", err)
	}

	return nil
}
