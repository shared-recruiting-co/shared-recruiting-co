package cloudfunctions

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"google.golang.org/api/gmail/v1"

	mail "github.com/shared-recruiting-co/libs/gmail"
)

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
	// TODO
	// 1. Fetch auth tokens for all user
	// 2. Spawn a goroutine for each user to watch their emails
	// 3. Wait for all goroutines to finish
	// 4. Mark success/failure in DB
	auth, err := jsonFromEnv("GOOGLE_AUTH_TOKEN")
	if err != nil {
		log.Fatalf("error getting auth token: %v", err)
		return
	}

	ctx := context.Background()
	srv, err := mail.NewGmailService(ctx, creds, auth)

	user := "me"
	label := "UNREAD"
	topic := os.Getenv("PUBSUB_TOPIC")

	// Watch for changes in labelId
	resp, err := srv.Users.Watch(user, &gmail.WatchRequest{
		LabelIds:  []string{label},
		TopicName: topic,
	}).Do()
	if err != nil {
		log.Fatalf("unable to watch: %v", err)
	}
	// success
	log.Println(resp)
}
