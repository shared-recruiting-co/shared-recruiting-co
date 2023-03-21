package cloudfunctions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/getsentry/sentry-go"

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
)

const (
	provider = "google"
	limit    = int32(1000)
)

func init() {
	functions.HTTP("CandidateGmailSubscription", candidateGmailSubscription)
	functions.HTTP("RecruiterGmailSubscription", recruiterGmailSubscription)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

// generic error handler
func handleError(w http.ResponseWriter, msg string, err error) {
	err = fmt.Errorf("%s: %w", msg, err)
	log.Print(err)
	sentry.CaptureException(err)
	w.WriteHeader(http.StatusInternalServerError)
}

type CloudFunction struct {
	ctx        context.Context
	oauthCreds []byte
	queries    db.Querier
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
		ctx:        ctx,
		oauthCreds: creds,
		queries:    queries,
	}, nil
}

type EmailSetting struct {
	IsActive bool `json:"is_active"`
}

type EmailSettings map[string]EmailSetting

func (cf *CloudFunction) isEmailActive(email string, inboxType db.InboxType) (bool, error) {
	if inboxType == db.InboxTypeCandidate {
		// get the user's profile
		profile, err := cf.queries.GetUserProfileByEmail(cf.ctx, email)
		if err != nil {
			return false, fmt.Errorf("error getting user profile: %w", err)
		}

		// check if the user's email is active
		return profile.IsActive, nil
	}

	if inboxType == db.InboxTypeRecruiter {
		// get the recruiter's profile
		profile, err := cf.queries.GetRecruiterByEmail(cf.ctx, email)
		if err != nil {
			return false, fmt.Errorf("error getting recruiter profile: %w", err)
		}

		emailSettings := EmailSettings{}
		err = json.Unmarshal([]byte(profile.EmailSettings), &emailSettings)
		if err != nil {
			return false, fmt.Errorf("error unmarshalling email settings: %w", err)
		}

		settings, ok := emailSettings[email]
		if !ok {
			return false, fmt.Errorf("email settings do not exists for: %s", email)
		}

		return settings.IsActive, nil
	}

	return false, fmt.Errorf("unsupported or invalid inbox type: %s", inboxType)
}

func (cf *CloudFunction) watch(users []db.UserOauthToken, arg *gmail.WatchRequest) []error {
	var err error
	var srv *srcmail.Service
	errs := []error{}

	for _, user := range users {
		auth := []byte(user.Token)

		srv, err = srcmail.NewService(cf.ctx, cf.oauthCreds, auth)
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
				err = cf.queries.UpsertUserOAuthToken(cf.ctx, db.UpsertUserOAuthTokenParams{
					UserID:   user.UserID,
					Email:    user.Email,
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
		userProfile, err := cf.queries.GetUserProfileByEmail(cf.ctx, gmailProfile.EmailAddress)
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
			return srv.Users.Watch(srv.UserID, arg).Do()
		})

		if err != nil {
			err = fmt.Errorf("error watching email: %w", err)
			errs = append(errs, err)
			continue
		}
	}

	return errs
}

func recruiterGmailSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("received watch trigger")
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

	topic := os.Getenv("PUBSUB_TOPIC")
	if topic == "" {
		log.Printf("missing PUBSUB_TOPIC")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	arg := &gmail.WatchRequest{
		LabelIds:          []string{"SENT"},
		LabelFilterAction: "include",
		TopicName:         topic,
	}

	hasError := false
	offset := int32(0)

	for {
		// 1. Fetch valid auth tokens
		userTokens, err := cf.queries.ListRecruiterOAuthTokens(ctx, db.ListRecruiterOAuthTokensParams{
			Provider: provider,
			IsValid:  true,
			Limit:    limit,
			Offset:   offset,
		})

		if err != nil {
			handleError(w, "error getting recruiter oauth tokens", err)
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

		errs := cf.watch(users, arg)

		for _, err := range errs {
			log.Printf("error watching email: %v", err)
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

func candidateGmailSubscription(w http.ResponseWriter, r *http.Request) {
	log.Println("received watch trigger")
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

	topic := os.Getenv("PUBSUB_TOPIC")
	if topic == "" {
		log.Printf("missing PUBSUB_TOPIC")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	arg := &gmail.WatchRequest{
		LabelIds:          []string{"UNREAD"},
		LabelFilterAction: "include",
		TopicName:         topic,
	}

	hasError := false
	offset := int32(0)

	for {
		// 1. Fetch valid auth tokens
		userTokens, err := cf.queries.ListCandidateOAuthTokens(ctx, db.ListCandidateOAuthTokensParams{
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

		errs := cf.watch(users, arg)

		for _, err := range errs {
			log.Printf("error watching email: %v", err)
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
