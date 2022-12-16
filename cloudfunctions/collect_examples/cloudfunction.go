package cloudfunctions

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/lib/pq"

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
)

const (
	provider = "google"
)

func init() {
	functions.HTTP("CollectExamples", collectExamples)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

// collectExamples collects examples from each user's inbox and forwards them to examples@sharedrecruiting.co
// In the future, we can do this in realtime whenever a email is labeled with @SRC, but for now we'll just do it ad-hoc in batch
func collectExamples(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		log.Printf("error fetching google app credentials: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 0, Create SRC client
	connectionURI := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		log.Printf("error connecting to database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	// use a max of 2 connections
	db.SetMaxOpenConns(2)

	// prepare queries
	queries, err := client.Prepare(ctx, db)
	if err != nil {
		log.Printf("error preparing db queries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 1. Fetch auth tokens for all user
	userTokens, err := queries.ListValidOAuthTokensByProvider(ctx, provider)
	if err != nil {
		log.Printf("error getting user tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hasError := false

	for _, userToken := range userTokens {
		log.Printf("processing user %s", userToken.UserID)
		// get the last scraped date from sync history
		history, err := queries.GetUserEmailSyncHistory(ctx, userToken.UserID)
		if err != nil {
			log.Printf("error getting user sync history: %v", err)
			hasError = true
			continue
		}

		auth := []byte(userToken.Token.RawMessage)

		srv, err := mail.NewService(ctx, creds, auth)
		if err != nil {
			log.Printf("error creating gmail service: %v", err)
			hasError = true
			continue
		}

		// Create recruiting detector client
		var messages []*gmail.Message
		pageToken := ""

		to := "Examples <examples@sharedrecruiting.co>"

		// batch process messages to reduce memory usage
		for {
			// Make Request to Fetch New Emails from Previous History ID
			// get next set of messages
			// if this is the first notification, trigger a one-time sync
			var startDate time.Time
			// valid is true if time is non-null
			if history.ExamplesCollectedAt.Valid {
				startDate = history.ExamplesCollectedAt.Time
			}
			// start
			messages, pageToken, err = fetchSRCEmails(srv, startDate, pageToken)

			// for now, abort on error
			if err != nil {
				log.Printf("error fetching emails: %v", err)
				hasError = true
				break
			}

			log.Printf("found %d messages to collect", len(messages))

			// forward each message
			for _, message := range messages {
				// payload isn't included in the list endpoint responses
				_, err := srv.ForwardEmail(message.Id, to)

				// for now, abort on error
				if err != nil {
					log.Printf("error sending message: %v", err)
					hasError = true
					continue
				}
			}

			if pageToken == "" {
				break
			}
		}

		// save sync date
		err = queries.UpsertUserEmailSyncHistory(ctx, client.UpsertUserEmailSyncHistoryParams{
			UserID:    userToken.UserID,
			HistoryID: history.HistoryID,
			SyncedAt:  history.SyncedAt,
			ExamplesCollectedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})

		if err != nil {
			log.Printf("error upserting user sync history: %v", err)
			hasError = true
			continue
		}
	}

	// write error status code for tracking
	if hasError {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("done.")
}
