package cloudfunctions

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
)

const (
	provider = "google"
)

func init() {
	functions.HTTP("FullEmailSync", fullEmailSync)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

type FullEmailSyncRequest struct {
	Email     string    `json:"email"`
	StartDate time.Time `json:"start_date"`
}

// fullEmailSync is triggers a sync up to the specified date for the given email
func fullEmailSync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	var data FullEmailSyncRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error unmarshalling request body: %v", err)
		http.Error(w, "Error unmarshalling request body", http.StatusInternalServerError)
		return
	}
	email := data.Email
	startDate := data.StartDate

	log.Println("full email sync request")

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

	// 1. Get User from email address
	user, err := queries.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("error getting user by email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, client.GetUserOAuthTokenParams{
		UserID:   user.ID,
		Provider: provider,
	})
	if err != nil {
		log.Printf("error getting user oauth token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create Gmail Service
	auth := []byte(userToken.Token.RawMessage)
	gmailSrv, err := mail.NewGmailService(ctx, creds, auth)
	gmailUser := "me"

	// Create SRC Labels
	_, err = mail.GetOrCreateSRCLabel(gmailSrv, gmailUser)
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
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
		log.Printf("error getting or creating the SRC label: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	srcJobOpportunityLabel, err := mail.GetOrCreateSRCJobOpportunityLabel(gmailSrv, gmailUser)
	if err != nil {
		log.Printf("error getting or creating the SRC job opportunity label: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create recruiting detector client
	classifierBaseURL := os.Getenv("CLASSIFIER_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, classifierBaseURL)
	if err != nil {
		log.Printf("error creating id token source: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		log.Printf("error getting id token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	classifier := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL:   classifierBaseURL,
		AuthToken: idToken.AccessToken,
	})

	var messages []*gmail.Message
	pageToken := ""

	// batch process messages to reduce memory usage
	for {
		// get next set of messages
		// if this is the first notification, trigger a one-time sync
		messages, pageToken, err = getEmailsSinceDate(gmailSrv, gmailUser, startDate, pageToken)

		// for now, abort on error
		if err != nil {
			log.Printf("error fetching emails: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// process messages
		examples := map[string]string{}
		for _, message := range messages {
			// payload isn't included in the list endpoint responses
			message, err := gmailSrv.Users.Messages.Get(gmailUser, message.Id).Do()

			// for now, abort on error
			if err != nil {
				log.Printf("error getting message: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// filter out empty messages
			if message.Payload == nil {
				continue
			}
			text, err := mail.MessageToString(message)
			examples[message.Id] = text
		}

		log.Printf("number of emails to classify: %d", len(examples))

		if len(examples) == 0 {
			break
		}

		// Batch predict on new emails
		results, err := classifier.PredictBatch(examples)
		if err != nil {
			log.Printf("error predicting on examples: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get IDs of new recruiting emails
		recruitingEmailIDs := []string{}
		for id, result := range results {
			if !result {
				continue
			}
			recruitingEmailIDs = append(recruitingEmailIDs, id)
		}

		log.Printf("number of recruiting emails: %d", len(recruitingEmailIDs))

		// Take action on recruiting emails
		if len(recruitingEmailIDs) > 0 {
			err = gmailSrv.Users.Messages.BatchModify(gmailUser, &gmail.BatchModifyMessagesRequest{
				Ids: recruitingEmailIDs,
				// Add SRC Job Label
				AddLabelIds: []string{srcJobOpportunityLabel.Id},
				// In future,
				// - mark as read
				// - archive
				// - create response
				// RemoveLabelIds: []string{"UNREAD"},
			}).Do()

			// for now, abort on error
			if err != nil {
				log.Printf("error modifying recruiting emails: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if pageToken == "" {
			break
		}
	}

	log.Printf("done.")
}
