package cloudfunctions

import (
	"encoding/base64"
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

const provider = "google"

func init() {
	functions.HTTP("RunWatchEmails", runWatchEmails)
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

func runWatchEmails(w http.ResponseWriter, r *http.Request) {
	log.Println("received watch trigger")
	ctx := r.Context()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       "watch-emails",
	})
	if err != nil {
		log.Printf("sentry.Init: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		handleError(w, "error getting credentials", err)
		return
	}

	// Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// TODO
	// v0 -> no pagination, no go routines
	// 2. Spawn a goroutine for each user to watch their emails
	// https://docs.sentry.io/platforms/go/concurrency/
	// 3. Wait for all goroutines to finish

	// 1. Fetch valid auth tokens for all users
	userTokens, err := queries.ListUserOAuthTokens(ctx, db.ListUserOAuthTokensParams{
		Provider: provider,
		IsValid:  true,
	})

	if err != nil {
		handleError(w, "error getting user tokens", err)
		return
	}

	var srv *srcmail.Service
	user := "me"
	label := "UNREAD"
	topic := os.Getenv("PUBSUB_TOPIC")

	hasError := false

	for _, userToken := range userTokens {
		auth := []byte(userToken.Token)

		srv, err = srcmail.NewService(ctx, creds, auth)
		if err != nil {
			err = fmt.Errorf("error creating gmail service: %w", err)
			log.Print(err)
			sentry.CaptureException(err)
			hasError = true
			continue
		}

		// Get the user's email address
		// This also keeps the user's refresh token valid for deactivated emails
		gmailProfile, err := srv.Profile()
		if err != nil {
			err = fmt.Errorf("error getting gmail profile: %w", err)
			log.Print(err)

			// check for oauth token expiration or revocation
			if srcmail.IsOAuth2Error(err) {
				log.Printf("error oauth error: %v", err)
				// update the user's oauth token
				err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
					UserID:   userToken.UserID,
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
			sentry.CaptureException(err)
			hasError = true
			continue
		}

		// validate the user's email is active
		userProfile, err := queries.GetUserProfileByEmail(ctx, gmailProfile.EmailAddress)
		if err != nil {
			err = fmt.Errorf("error getting user profile: %w", err)
			log.Print(err)
			sentry.CaptureException(err)
			hasError = true
			continue
		}

		if !userProfile.IsActive {
			log.Printf("skipping deactivated email %s", userProfile.Email)
			continue
		}

		// Watch for changes in labelId
		resp, err := srcmail.ExecuteWithRetries(func() (*gmail.WatchResponse, error) {
			return srv.Users.Watch(user, &gmail.WatchRequest{
				LabelIds:          []string{label},
				LabelFilterAction: "include",
				TopicName:         topic,
			}).Do()
		})

		if err != nil {
			err = fmt.Errorf("error watching email: %w", err)
			log.Print(err)
			sentry.CaptureException(err)
			continue
		}
		// success
		log.Printf("watching: %v", resp)
	}

	// write error status code for tracking
	if hasError {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("done.")
}
