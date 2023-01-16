package cloudfunctions

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/gmail/v1"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
)

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
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Fetch valid auth tokens for all users
	userTokens, err := queries.ListUserOAuthTokens(ctx, db.ListUserOAuthTokensParams{
		Provider: provider,
		IsValid:  true,
	})

	if err != nil {
		log.Printf("error getting user tokens: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var srv *srcmail.Service
	hasError := false

	for _, userToken := range userTokens {
		auth := []byte(userToken.Token)
		srv, err = srcmail.NewService(ctx, creds, auth)
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
		oldLabel, err := srv.GetOrCreateLabel(&gmail.Label{
			Name:  "@SRC/Job Opportunity",
			Color: srclabel.JobsOpportunity.Color,
		})
		if err != nil {
			log.Printf("error getting old labels: %v", err)
			hasError = true
			continue
		}

		// Get all messages with the old label scheme and not new scheme
		q := fmt.Sprintf("-label:%s label:%s", srclabel.JobsOpportunity.Name, oldLabel.Name)
		pageToken := ""
		fetchMessageError := false
		for {
			messages, err := srv.Users.Messages.List(srv.UserID).Q(q).PageToken(pageToken).MaxResults(500).Do()
			if err != nil {
				log.Printf("error getting messages: %v", err)
				hasError = true
				fetchMessageError = true
				break
			}

			messageIDs := make([]string, len(messages.Messages))
			for i, message := range messages.Messages {
				messageIDs[i] = message.Id
			}

			if len(messageIDs) > 0 {
				// Remove old labels
				err = srv.Users.Messages.BatchModify(srv.UserID, &gmail.BatchModifyMessagesRequest{
					Ids:            messageIDs,
					RemoveLabelIds: []string{newLabels.SRC.Id, oldLabel.Id},
				}).Do()

				if err != nil {
					log.Printf("error removing old labels: %v", err)
					hasError = true
					fetchMessageError = true
					break
				}
			}

			if messages.NextPageToken == "" {
				break
			}
			pageToken = messages.NextPageToken
		}

		if fetchMessageError {
			continue
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
func syncLabels(srv *srcmail.Service, labels *srclabel.Labels) error {
	// Update each label to update properties
	_, err := srv.Users.Labels.Update(srv.UserID, labels.SRC.Id, &srclabel.SRC).Do()
	if err != nil {
		return err
	}

	// The rest of the labels are new, so we need to update them
	return nil
}
