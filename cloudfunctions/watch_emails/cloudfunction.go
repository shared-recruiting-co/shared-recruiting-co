package cloudfunctions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/lib/pq"

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
		log.Fatalf("error getting credentials: %v", err)
		return
	}

	// Create SRC client
	ctx := context.Background()

	connectionURI := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionURI)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
		return
	}

	queries := client.New(db)

	// TODO
	// v0 -> no pagination, no go routines
	// 2. Spawn a goroutine for each user to watch their emails
	// 3. Wait for all goroutines to finish
	// 4. Mark success/failure in DB

	// 1. Fetch auth tokens for all user
	userTokens, err := queries.ListOAuthTokensByProvider(ctx, provider)
	if err != nil {
		log.Fatalf("error getting user tokens: %v", err)
		return
	}

	var srv *gmail.Service
	user := "me"
	label := "UNREAD"
	topic := os.Getenv("PUBSUB_TOPIC")

	for _, userToken := range userTokens {
		auth := []byte(userToken.Token.RawMessage)

		srv, err = mail.NewGmailService(ctx, creds, auth)
		if err != nil {
			log.Fatalf("error creating gmail service: %v", err)
		}
		// Watch for changes in labelId
		resp, err := srv.Users.Watch(user, &gmail.WatchRequest{
			LabelIds:  []string{label},
			TopicName: topic,
		}).Do()
		if err != nil {
			log.Fatalf("error watching: %v", err)
		}
		// success
		log.Println(resp)
	}
	log.Println("done.")
}
