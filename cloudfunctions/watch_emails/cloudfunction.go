package cloudfunctions

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
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

func runWatchEmails(w http.ResponseWriter, r *http.Request) {
	log.Println("received watch trigger")
	creds, err := jsonFromEnv("GOOGLE_APPLICATION_CREDENTIALS")
	if err != nil {
		log.Printf("error getting credentials: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create SRC client
	ctx := r.Context()

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

	queries := client.New(db)

	// TODO
	// v0 -> no pagination, no go routines
	// 2. Spawn a goroutine for each user to watch their emails
	// 3. Wait for all goroutines to finish

	// 1. Fetch auth tokens for all user
	userTokens, err := queries.ListValidOAuthTokensByProvider(ctx, provider)
	if err != nil {
		log.Printf("error getting user tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var srv *gmail.Service
	user := "me"
	label := "UNREAD"
	topic := os.Getenv("PUBSUB_TOPIC")

	hasError := false

	for _, userToken := range userTokens {
		auth := []byte(userToken.Token.RawMessage)

		srv, err = mail.NewGmailService(ctx, creds, auth)
		if err != nil {
			log.Printf("error creating gmail service: %v", err)
			hasError = true
			continue
		}
		// Watch for changes in labelId
		resp, err := srv.Users.Watch(user, &gmail.WatchRequest{
			LabelIds:  []string{label},
			TopicName: topic,
		}).Do()

		if err != nil {
			log.Printf("error watching: %v", err)
			// check for oauth token expiration or revocation
			oauth2Err := &oauth2.RetrieveError{}
			if errors.As(err, &oauth2Err) {
				log.Printf("error oauth error: %v", oauth2Err)
				// update the user's oauth token
				err = queries.UpsertUserOAuthToken(ctx, client.UpsertUserOAuthTokenParams{
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
			hasError = true
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
