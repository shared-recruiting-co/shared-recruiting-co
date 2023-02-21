package cloudfunctions

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
)

// reclassify is triggers a sync up to the specified date for the given email
func reclassify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       "reclassify",
	})
	if err != nil {
		log.Printf("sentry.Init: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	log.Println("reclassify request")

	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		handleError(w, "error fetching google app credentials", err)
		return
	}

	// Create recruiting detector client
	classifierBaseURL := os.Getenv("ML_SERVICE_URL")
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

	classifier := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   classifierBaseURL,
		AuthToken: idToken.AccessToken,
	})

	// 0, Create SRC http client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	userTokens, err := queries.ListUserOAuthTokens(ctx, db.ListUserOAuthTokensParams{
		Provider: "google",
		IsValid:  true,
	})
	if err != nil {
		handleError(w, "error fetching user tokens", err)
		return
	}

	for _, userToken := range userTokens {
		// Get User' OAuth Token
		userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
			UserID:   userToken.UserID,
			Provider: provider,
		})
		if err != nil {
			handleError(w, "error getting user oauth token", err)
			return
		}

		// Create Gmail Service
		auth := []byte(userToken.Token)
		srv, err := srcmail.NewService(ctx, creds, auth)
		if err != nil {
			handleError(w, "error creating mail service", err)
			return
		}

		// Create SRC Labels
		labels, err := srv.GetOrCreateCandidateLabels()
		if err != nil {
			// first request, so check if the error is an oauth error
			// if so, update the database
			if srcmail.IsOAuth2Error(err) {
				log.Printf("error oauth error: %v", err)
				// update the user's oauth token
				err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
					UserID:   userToken.UserID,
					Email:    userToken.Email,
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

		var threads []*gmail.Thread
		pageToken := ""

		// batch process messages to reduce memory usage
		for {
			// get next set of messages
			// if this is the first notification, trigger a one-time sync
			threads, pageToken, err = fetchRecruitingThreads(srv, pageToken)

			// for now, abort on error
			if err != nil {
				handleError(w, "error fetching emails", err)
				return
			}

			// get the messages for each thread
			var messages []*gmail.Message
			for _, t := range threads {
				thread, err := srv.GetThread(t.Id, "minimal")
				if err != nil {
					// for now abort on error
					handleError(w, "error fetching thread", err)
					return
				}

				// get messages before the first reply
				filtered := filterMessagesAfterReply(thread.Messages)
				// save for processing
				messages = append(messages, filtered...)
			}

			// process messages
			examples := map[string]*ml.ClassifyRequest{}
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
				examples[message.Id] = &ml.ClassifyRequest{
					From:    srcmessage.Sender(message),
					Subject: srcmessage.Subject(message),
					Body:    srcmessage.Body(message),
				}
			}

			log.Printf("number of emails to re-classify: %d", len(examples))

			if len(examples) == 0 {
				break
			}

			// Batch predict on new emails
			results, err := classifier.BatchClassify(&ml.BatchClassifyRequest{
				Inputs: examples,
			})
			if err != nil {
				handleError(w, "error predicting on examples", err)
				return
			}

			// Get IDs of messages that are NO LONGER classified as recruiting
			nonRecruitingIDs := []string{}
			for id, result := range results.Results {
				if result {
					continue
				}
				nonRecruitingIDs = append(nonRecruitingIDs, id)
			}

			log.Printf("number of non-recruiting emails: %d", len(nonRecruitingIDs))

			// Take action on recruiting emails
			err = handleNonRecruitingEmails(srv, labels, nonRecruitingIDs)

			// for now, abort on error
			if err != nil {
				handleError(w, "error modifying recruiting emails", err)
				return
			}

			if pageToken == "" {
				break
			}
		}
	}

	log.Printf("done.")
}

func handleNonRecruitingEmails(srv *srcmail.Service, labels *srclabel.CandidateLabels, messageIDs []string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	_, err := srcmail.ExecuteWithRetries(func() (interface{}, error) {
		err := srv.Users.Messages.BatchModify(srv.UserID, &gmail.BatchModifyMessagesRequest{
			Ids: messageIDs,
			// Remove job opportunity label and parent folder labels
			RemoveLabelIds: []string{labels.SRC.Id, labels.Jobs.Id, labels.JobsOpportunity.Id},
		}).Do()
		// hack to make compatible with ExecuteWithRetries requirements
		return nil, err
	})

	if err != nil {
		return fmt.Errorf("error modifying recruiting emails: %v", err)
	}

	return nil
}

// fetchRecruitingThreads fetches all threads since the start date
// It ignores threads of only sent emails and threads already processed by SRC
func fetchRecruitingThreads(srv *srcmail.Service, pageToken string) ([]*gmail.Thread, string, error) {
	// get all existing job threads
	q := fmt.Sprintf("label:%s", srclabel.JobsOpportunity.Name)

	r, err := srcmail.ExecuteWithRetries(func() (*gmail.ListThreadsResponse, error) {
		return srv.Users.Threads.
			List(srv.UserID).
			PageToken(pageToken).
			Q(q).
			MaxResults(250).
			Do()
	})

	if err != nil {
		return nil, "", err
	}

	return r.Threads, r.NextPageToken, nil
}
