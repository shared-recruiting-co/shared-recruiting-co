package cloudfunctions

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
)

const (
	provider                = "google"
	SRC_Label               = "@SRC"
	SRC_JobOpportunityLabel = "@SRC/Job Opportunity"
	SRC_Color               = "#16a765"
	White                   = "#ffffff"
)

func init() {
	functions.HTTP("NewUserWorkflow", newUserWorkflow)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

// newUserWorkflow
func newUserWorkflow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get the authorization header from the request
	authHeader := r.Header.Get("Authorization")
	// if the authorization header is empty, return an error
	if authHeader == "" {
		log.Printf("missing authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// get the token from the authorization header
	tokenRaw := authHeader[len("Bearer "):]
	// get secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Printf("missing JWT_SECRET env var")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(tokenRaw, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("error parsing jwt token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		log.Printf("invalid jwt token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	// if the token is invalid, return an error
	if !ok {
		log.Print("error getting jwt claims")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get the user id from the subject of the token
	sub := claims["sub"].(string)

	// if the user id is empty, return an error
	if sub == "" {
		log.Printf("user id (sub) is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(sub)

	if err != nil {
		log.Printf("error parsing user id into uuid: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	creds, err := jsonFromEnv("GOOGLE_APPLICATION_CREDENTIALS")
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

	// prepare queries
	queries, err := client.Prepare(ctx, db)
	if err != nil {
		log.Printf("error preparing db queries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, client.GetUserOAuthTokenParams{
		UserID:   userID,
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
	srcLabel, err := mail.GetOrCreateLabel(gmailSrv, gmailUser, SRC_Label, SRC_Color, White)
	if err != nil {
		log.Printf("error getting or creating the SRC label: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	srcJobOpportunityLabel, err := mail.GetOrCreateLabel(gmailSrv, gmailUser, SRC_JobOpportunityLabel, SRC_Color, White)
	if err != nil {
		log.Printf("error getting or creating the SRC job opportunity label: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create recruiting detector client
	classifier := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL: os.Getenv("CLASSIFIER_URL"),
		ApiKey:  os.Getenv("CLASSIFIER_API_KEY"),
	})

	var messages []*gmail.Message
	pageToken := ""

	// batch process messages to reduce memory usage
	for {
		// Make Request to Fetch New Emails from Previous History ID
		// get next set of messages
		// if this is the first notification, trigger a one-time sync
		messages, pageToken, err = GetEmailsSinceLastYear(gmailSrv, gmailUser, pageToken)

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
				// Add SRC Label
				AddLabelIds: []string{srcLabel.Id, srcJobOpportunityLabel.Id},
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

	// Watch for changes in unread messages
	topic := os.Getenv("PUBSUB_TOPIC")
	resp, err := gmailSrv.Users.Watch(gmailUser, &gmail.WatchRequest{
		LabelIds:  []string{"UNREAD"},
		TopicName: topic,
	}).Do()

	if err != nil {
		log.Printf("error watching for unread messages: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("watch response: %v", resp)

	err = queries.UspertUserEmailSyncHistory(ctx, client.UspertUserEmailSyncHistoryParams{
		UserID:    userID,
		HistoryID: int64(resp.HistoryId),
	})

	if err != nil {
		log.Printf("error upserting user email sync history: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("done.")
}
