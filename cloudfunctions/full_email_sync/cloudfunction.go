package cloudfunctions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/getsentry/sentry-go"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
	mail "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/gmail/label"
)

const (
	provider = "google"
)

var (
	// global variable to share across functions...simplest approach for now
	examplesCollectorSrv   *mail.Service
	collectedExampleLabels = []string{"INBOX", "UNREAD"}
)

func init() {
	functions.HTTP("FullEmailSync", fullEmailSync)
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
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type FullEmailSyncRequest struct {
	Email     string    `json:"email"`
	StartDate time.Time `json:"start_date"`
}

// fullEmailSync is triggers a sync up to the specified date for the given email
func fullEmailSync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Printf("sentry.Init: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	// Get the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, "error reading request body", err)
		return
	}
	var data FullEmailSyncRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		handleError(w, "error unmarshalling request body", err)
		return
	}
	email := data.Email
	startDate := data.StartDate

	log.Println("full email sync request")

	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		handleError(w, "error fetching google app credentials", err)
		return
	}

	// 0, Create SRC http client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := client.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetUserProfileByEmail(ctx, email)
	if err != nil {
		handleError(w, "error getting user profile by email", err)
		return
	}
	// if auto contribute is on, create the collector service
	if user.AutoContribute {
		auth, err := jsonFromEnv("EXAMPLES_GMAIL_OAUTH_TOKEN")
		if err != nil {
			handleError(w, "error reading examples@sharedrecruiting.co credentials", err)
			return
		}
		examplesCollectorSrv, err = mail.NewService(ctx, creds, auth)
		if err != nil {
			handleError(w, "error creating example collector service", err)
			return
		}
	}

	// Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, client.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Provider: provider,
	})
	if err != nil {
		handleError(w, "error getting user oauth token", err)
		return
	}

	// Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := mail.NewService(ctx, creds, auth)

	// Create SRC Labels
	labels, err := srv.GetOrCreateSRCLabels()
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
		if mail.IsOAuth2Error(err) {
			log.Printf("error oauth error: %v", err)
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
		handleError(w, "error getting or creating the SRC label", err)
		return
	}

	// Create recruiting detector client
	classifierBaseURL := os.Getenv("CLASSIFIER_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, classifierBaseURL)
	if err != nil {
		handleError(w, "error creating id token source", err)
		return
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		handleError(w, "error getting id token", err)
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
		messages, pageToken, err = fetchEmailsSinceDate(srv, startDate, pageToken)

		// for now, abort on error
		if err != nil {
			handleError(w, "error fetching emails", err)
			return
		}

		// process messages
		examples := map[string]*PredictRequest{}
		for _, message := range messages {
			// payload isn't included in the list endpoint responses
			message, err := srv.GetMessage(message.Id)

			// for now, abort on error
			if err != nil {
				handleError(w, "error getting message", err)
				return
			}

			// filter out empty messages
			if message.Payload == nil {
				continue
			}
			example := &PredictRequest{
				From:    mail.MessageSender(message),
				Subject: mail.MessageSubject(message),
				Body:    mail.MessageBody(message),
			}
			examples[message.Id] = example
		}

		log.Printf("number of emails to classify: %d", len(examples))

		if len(examples) == 0 {
			break
		}

		// Batch predict on new emails
		results, err := classifier.PredictBatch(examples)
		if err != nil {
			handleError(w, "error predicting on examples", err)
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
		err = handleRecruitingEmails(srv, user, labels, recruitingEmailIDs)
		// for now, abort on error
		if err != nil {
			handleError(w, "error modifying recruiting emails", err)
			return
		}

		// save statistics
		if len(examples) > 0 {
			err = queries.IncrementUserEmailStat(
				ctx,
				client.IncrementUserEmailStatParams{
					UserID:    user.UserID,
					Email:     user.Email,
					StatID:    "emails_processed",
					StatValue: int32(len(examples)),
				},
			)
			if err != nil {
				// print error, but don't abort
				err = fmt.Errorf("error incrementing user email stat: %w", err)
				log.Print(err)
				sentry.CaptureException(err)
			}
		}
		if len(recruitingEmailIDs) > 0 {
			err = queries.IncrementUserEmailStat(
				ctx,
				client.IncrementUserEmailStatParams{
					UserID:    user.UserID,
					Email:     user.Email,
					StatID:    "jobs_detected",
					StatValue: int32(len(recruitingEmailIDs)),
				},
			)
			if err != nil {
				// print error, but don't abort
				err = fmt.Errorf("error incrementing user email stat: %w", err)
				log.Print(err)
				sentry.CaptureException(err)
			}
		}

		if pageToken == "" {
			break
		}
	}

	log.Printf("done.")
}

func handleRecruitingEmails(srv *mail.Service, profile client.UserProfile, labels *srclabel.Labels, messageIDs []string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	removeLabels := []string{}
	if profile.AutoArchive {
		removeLabels = append(removeLabels, "INBOX", "UNREAD")
	}

	_, err := mail.ExecuteWithRetries(func() (interface{}, error) {
		err := srv.Users.Messages.BatchModify(srv.UserID, &gmail.BatchModifyMessagesRequest{
			Ids: messageIDs,
			// Add job opportunity label and parent folder labels
			AddLabelIds:    []string{labels.SRC.Id, labels.Jobs.Id, labels.JobsOpportunity.Id},
			RemoveLabelIds: removeLabels,
		}).Do()
		// hack to make compatible with ExecuteWithRetries requirements
		return nil, err
	})

	if err != nil {
		return fmt.Errorf("error modifying recruiting emails: %v", err)
	}

	if profile.AutoContribute {
		for _, id := range messageIDs {
			// shouldn't happen
			if examplesCollectorSrv == nil {
				log.Print("examples collector service not initialized")
				break
			}
			// clone the message to the examples inbox
			_, err := mail.CloneMessage(srv, examplesCollectorSrv, id, collectedExampleLabels)

			if err != nil {
				// don't abort on error
				log.Printf("error collecting email %s: %v", id, err)
				continue
			}
		}
	}

	return nil
}
