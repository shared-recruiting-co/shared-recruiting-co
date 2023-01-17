package cloudfunctions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/getsentry/sentry-go"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
)

const provider = "google"

func init() {
	functions.HTTP("PopulateJobs", populateJobs)
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
	w.WriteHeader(http.StatusInternalServerError)
}

func filterMessagesAfterReply(messages []*gmail.Message) []*gmail.Message {
	filtered := []*gmail.Message{}
	// ensure messages are sorted by ascending date
	srcmessage.SortByDate(messages)

	for _, m := range messages {
		if srcmessage.IsSent(m) {
			break
		}
		filtered = append(filtered, m)
	}
	return filtered
}

// --------------------------------
// Daily Parse Job
// For each user
// // List their jobs
// // Get the latest job emailed_at date and thread ID
// // Fetch all job threads since that date
// // Filter out threads that have already been parsed
// // For each thread
// // // Get the message ID with the label
// // // Parse the message
// // // Save the job
// --------------------------------

// Populate jobs parses all unprocessed job threads and saves them to the database
func populateJobs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       "populate-jobs",
	})
	if err != nil {
		log.Printf("sentry.Init: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	log.Println("populate jobs request")

	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		handleError(w, "error fetching google app credentials", err)
		return
	}

	// Create recruiting detector client
	mlServiceBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, mlServiceBaseURL)
	if err != nil {
		handleError(w, "error creating id token source", err)
		return
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		handleError(w, "error getting id token", err)
		return
	}

	parser := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   mlServiceBaseURL,
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

		profile, err := srv.Profile()
		if err != nil {
			if srcmail.IsOAuth2Error(err) {
				log.Printf("error oauth error: %v", err)
				// update the user's oauth token
				err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
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
			handleError(w, "error getting user gmail profile", err)
			return
		}

		// 1. Get the current user's jobs
		existingJobs, err := queries.ListUserEmailJobs(ctx, db.ListUserEmailJobsParams{
			UserID: userToken.UserID,
			// assume a user isn't receiving more than 10 jobs a day
			Limit: 10,
			// get the first page of results (sorted desc by date received)
			Offset: 0,
		})
		if err != nil {
			handleError(w, "error fetching existing user jobs", err)
			return
		}

		// handle the case where the user has no jobs
		if len(existingJobs) == 0 {
			log.Printf("user %s has no jobs", userToken.UserID)
		}

		latestJobDate := time.Time{}

		// 2. Get the latest job's date and thread ID
		if len(existingJobs) > 0 {
			latestJobDate = existingJobs[0].EmailedAt
		}
		log.Printf("fetching threads since %s", latestJobDate.Format(time.RFC3339))

		var threads []*gmail.Thread
		pageToken := ""

		// batch process messages to reduce memory usage
		for {
			// get next set of messages
			// if this is the first notification, trigger a one-time sync
			threads, pageToken, err = fetchRecruitingThreads(srv, latestJobDate, pageToken)

			// for now, abort on error
			if err != nil {
				handleError(w, "error fetching emails", err)
				return
			}

			// get the messages for each thread
			messages := map[string]*gmail.Message{}

			for _, t := range threads {
				// avoid parsing the same message twice
				// in theory, this shouldn't happen because of latestJobDate
				if containsThread(existingJobs, t.Id) {
					log.Printf("skipping thread %s because it already exists", t.Id)
					continue
				}

				thread, err := srv.GetThread(t.Id, "minimal")
				if err != nil {
					// for now abort on error
					handleError(w, "error fetching thread", err)
					return
				}
				// get messages before the first reply
				filtered := filterMessagesAfterReply(thread.Messages)

				if len(filtered) == 0 {
					// no messages before the first reply
					continue
				}
				// get the earliest message (filterMessagesAfterReply sorts by date ascending)
				earliest := filtered[0]

				// save for processing
				messages[earliest.Id] = earliest
			}

			// process messages
			// process messages
			inputs := map[string]*ml.ParseJobRequest{}
			for id := range messages {
				// payload isn't included in the list endpoint responses
				message, err := srv.GetMessage(id)

				// for now, abort on error
				if err != nil {
					handleError(w, "error getting message", err)
					return
				}

				messages[id] = message

				// filter out empty messages
				if message.Payload == nil {
					continue
				}

				inputs[message.Id] = &ml.ParseJobRequest{
					From:    srcmessage.Sender(message),
					Subject: srcmessage.Subject(message),
					Body:    srcmessage.Body(message),
				}
			}

			log.Printf("number of emails to parse: %d", len(inputs))

			if len(inputs) == 0 {
				break
			}

			if len(inputs) == 0 {
				break
			}

			// predict one at a time for now
			for id, input := range inputs {
				log.Printf("parsing email: %s", id)
				job, err := parser.ParseJob(input)
				// for now, abort on error
				if err != nil {
					handleError(w, "error parsing job", err)
					return
				}

				// if fields are missing, skip
				if job.Company == "" || job.Title == "" || job.Recruiter == "" {
					// print sender and subject
					log.Printf("skipping job: %v", job)
					continue
				}

				message := messages[id]
				recruiterEmail := srcmessage.SenderEmail(message)
				data := map[string]interface{}{
					"recruiter":       job.Recruiter,
					"recruiter_email": recruiterEmail,
				}

				// turn data into json.RawMessage
				b, err := json.Marshal(data)
				if err != nil {
					handleError(w, "error marshaling data", err)
					return
				}

				// convert epoch ms to time.Time
				emailedAt := time.Unix(message.InternalDate/1000, 0)

				err = queries.InsertUserEmailJob(ctx, db.InsertUserEmailJobParams{
					UserID:        userToken.UserID,
					UserEmail:     profile.EmailAddress,
					EmailThreadID: message.ThreadId,
					EmailedAt:     emailedAt,
					Company:       job.Company,
					JobTitle:      job.Title,
					Data:          b,
				})

				// for now, continue on error
				if err != nil {
					log.Printf("error inserting job (%v): %v", job, err)
					continue
				}
			}

			if pageToken == "" {
				break
			}
		}
	}

	log.Printf("done.")
}

// fetchRecruitingThreads fetches all threads since the start date
// It ignores threads of only sent emails and threads already processed by SRC
func fetchRecruitingThreads(srv *srcmail.Service, startDate time.Time, pageToken string) ([]*gmail.Thread, string, error) {
	// get all existing job threads
	q := fmt.Sprintf("label:%s", srclabel.JobsOpportunity.Name)
	if !startDate.IsZero() {
		// use Unix time (seconds) to avoid timezone issues
		// Gmail InternalDate is in Unix
		q = fmt.Sprintf("%s after:%d", q, startDate.Unix())
	}

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

func containsThread(jobs []db.UserEmailJob, threadID string) bool {
	for _, j := range jobs {
		if j.EmailThreadID == threadID {
			return true
		}
	}
	return false
}
