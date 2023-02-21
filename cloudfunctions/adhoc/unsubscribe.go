package cloudfunctions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"gopkg.in/guregu/null.v4"

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
)

const (
	limit = int32(1000)
)

type CloudFunction struct {
	ctx                 context.Context
	OAuthCredententials []byte
	Queries             db.Querier
}

func NewCloudFunction(ctx context.Context) (*CloudFunction, error) {
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("error getting credentials: %w", err)
	}

	// Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("missing SUPABASE_API_URL")
	}
	apiKey := os.Getenv("SUPABASE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing SUPABASE_API_KEY")
	}
	queries := db.NewHTTP(apiURL, apiKey)
	// set debug logging
	queries.Debug = true

	return &CloudFunction{
		ctx:                 ctx,
		OAuthCredententials: creds,
		Queries:             queries,
	}, nil
}

func (cf *CloudFunction) stop(users []db.UserOauthToken) []error {
	var err error
	var srv *srcmail.Service
	errs := []error{}

	for _, user := range users {
		auth := []byte(user.Token)

		srv, err = srcmail.NewService(cf.ctx, cf.OAuthCredententials, auth)
		if err != nil {
			err = fmt.Errorf("error creating gmail service: %w", err)
			errs = append(errs, err)
			continue
		}

		// Get the user's email address
		// This also keeps the user's refresh token valid for deactivated emails
		gmailProfile, err := srv.Profile()
		if err != nil {
			err = fmt.Errorf("error getting gmail profile: %w", err)
			errs = append(errs, err)

			// check for oauth token expiration or revocation
			if srcmail.IsOAuth2Error(err) {
				log.Printf("error oauth error: %v", err)
				// update the user's oauth token
				err = cf.Queries.UpsertUserOAuthToken(cf.ctx, db.UpsertUserOAuthTokenParams{
					UserID:   user.UserID,
					Email:    null.StringFrom(gmailProfile.EmailAddress),
					Provider: provider,
					Token:    user.Token,
					IsValid:  false,
				})
				if err != nil {
					log.Printf("error updating user oauth token: %v", err)
				} else {
					log.Printf("marked user oauth token as invalid")
				}
			}
			sentry.CaptureException(err)
			continue
		}

		// validate the user's email is active
		userProfile, err := cf.Queries.GetUserProfileByEmail(cf.ctx, gmailProfile.EmailAddress)
		if err != nil {
			err = fmt.Errorf("error getting user profile: %w", err)
			errs = append(errs, err)
			continue
		}

		if !userProfile.IsActive {
			log.Printf("skipping deactivated email %s", userProfile.Email)
			continue
		}

		_, err = srcmail.ExecuteWithRetries(func() (*gmail.WatchResponse, error) {
			err := srv.Users.Stop(srv.UserID).Do()
			return nil, err
		})

		if err != nil {
			err = fmt.Errorf("error unsubscribing from email: %w", err)
			errs = append(errs, err)
			continue
		}
	}

	return errs
}

func unsubscribe(w http.ResponseWriter, r *http.Request) {
	log.Println("received unsubscribe trigger")
	ctx := r.Context()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       os.Getenv("FUNCTION_NAME"),
	})
	if err != nil {
		log.Printf("sentry.Init: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	cf, err := NewCloudFunction(ctx)
	if err != nil {
		log.Printf("error creating cloud function: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hasError := false
	offset := int32(0)

	for {
		// 1. Fetch valid auth tokens
		userTokens, err := cf.Queries.ListCandidateOAuthTokens(ctx, db.ListCandidateOAuthTokensParams{
			Provider: provider,
			IsValid:  true,
			Limit:    limit,
			Offset:   offset,
		})

		if err != nil {
			handleError(w, "error getting candidate oauth tokens", err)
			return
		}

		// convert to user oauth token
		users := make([]db.UserOauthToken, len(userTokens))
		for i, userToken := range userTokens {
			users[i] = db.UserOauthToken{
				UserID:   userToken.UserID,
				Provider: userToken.Provider,
				Token:    userToken.Token,
			}
		}

		errs := cf.stop(users)

		for _, err := range errs {
			log.Printf("error unsubscribing email: %v", err)
			sentry.CaptureException(err)
			hasError = true
		}

		// check if there are more results
		if len(userTokens) < int(limit) {
			break
		}
		offset += limit
	}

	if hasError {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("done.")
}
