package cloudfunctions

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail/label"
)

const provider = "google"

func init() {
	functions.HTTP("MigrateLabels", migrateLabels)
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

func migrateLabels(w http.ResponseWriter, r *http.Request) {
	log.Println("running labels migration")
	ctx := r.Context()
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		log.Printf("error getting credentials: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := client.NewHTTP(apiURL, apiKey)

	// 1. Fetch valid auth tokens for all users
	userTokens, err := queries.ListUserOAuthTokens(ctx, client.ListUserOAuthTokensParams{
		Provider: provider,
		IsValid:  true,
	})

	if err != nil {
		log.Printf("error getting user tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var srv *mail.Service
	hasError := false

	for _, userToken := range userTokens {
		auth := []byte(userToken.Token)
		srv, err = mail.NewService(ctx, creds, auth)
		if err != nil {
			log.Printf("error creating gmail service: %v", err)
			hasError = true
			continue
		}
		// Get new labels
		newLabels, err := srv.GetOrCreateSRCLabels()
		if err != nil {
			log.Printf("error getting new labels: %v", err)
			hasError = true
			continue
		}

		// Get old label Ids
		oldLabel, err := srv.GetOrCreateSRCJobOpportunityLabel()
		if err != nil {
			log.Printf("error getting old labels: %v", err)
			hasError = true
			continue
		}

		// Get all messages with the old label scheme and not new scheme
		q := fmt.Sprintf("label:%s label:%s -label:%s", srclabel.SRC.Name, srclabel.JobsOpportunity.Name, oldLabel.Name)
		messages, err := srv.Users.Messages.List(srv.UserID).Q(q).Do()
		if err != nil {
			log.Printf("error getting messages: %v", err)
			hasError = true
			continue
		}

		messageIDs := make([]string, len(messages.Messages))
		for i, message := range messages.Messages {
			messageIDs[i] = message.Id
		}

		if len(messageIDs) > 0 {
			// Remove old labels
			err = srv.Users.Messages.BatchModify(srv.UserID, &gmail.BatchModifyMessagesRequest{
				Ids: messageIDs,
				// Add job opportunity label and parent folder labels
				RemoveLabelIds: []string{newLabels.SRC.Id, oldLabel.Id},
			}).Do()

			if err != nil {
				log.Printf("error removing old labels: %v", err)
				hasError = true
				continue
			}
		}

		// Delete old job label
		err = srv.Users.Labels.Delete(srv.UserID, oldLabel.Id).Do()
		if err != nil {
			log.Printf("error deleting old label: %v", err)
			hasError = true
			continue
		}

		// Sync new label
		err = syncLabels(srv, newLabels)
		if err != nil {
			log.Printf("error syncing labels: %v", err)
			hasError = true
			continue
		}
		log.Printf("successfully migrated labels for user %s", userToken.UserID)
	}

	// write error status code for tracking
	if hasError {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("done.")
}

// Consider moving into shared lib if this becomes a common operation
func syncLabels(srv *mail.Service, labels *srclabel.Labels) error {
	// Update each label to update properties
	_, err := srv.Users.Labels.Update(srv.UserID, labels.SRC.Id, &srclabel.SRC).Do()
	if err != nil {
		return err
	}

	// The rest of the labels are new, so we need to update them
	return nil
}
