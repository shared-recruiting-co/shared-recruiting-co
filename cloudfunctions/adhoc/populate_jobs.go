package cloudfunctions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
)

// 1. Fetch all threads with srclabel.JobsOpportunity
// 2. Filter all messages before a reply
// 3. Pick the earliest message
// 4. Hit the parse endpoint
// 5. if all fields are present, insert a job
// 6. if any fields are missing, log

type PopulateJobsRequest struct {
	Email     string    `json:"email"`
	StartDate time.Time `json:"start_date"`
}

func handleError(w http.ResponseWriter, msg string, err error) {
	err = fmt.Errorf("%s: %w", msg, err)
	log.Print(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func populateJobs(w http.ResponseWriter, r *http.Request) {
	log.Println("running populate jobs")
	ctx := r.Context()
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		log.Printf("error getting credentials: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, "error reading request body", err)
		return
	}
	var data PopulateJobsRequest
	err = json.Unmarshal(body, &data)
	if err != nil {
		handleError(w, "error unmarshalling request body", err)
		return
	}
	email := data.Email
	startDate := data.StartDate

	log.Println("populate jobs request")

	// Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Fetch valid auth tokens for all users
	user, err := queries.GetUserProfileByEmail(ctx, email)
	if err != nil {
		handleError(w, "error getting user profile by email", err)
		return
	}

	// Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Provider: provider,
	})
	if err != nil {
		handleError(w, "error getting user oauth token", err)
		return
	}

	// Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := srcmail.NewService(ctx, creds, auth)

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

	mlSrv := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   mlServiceBaseURL,
		AuthToken: idToken.AccessToken,
	})

	var threads []*gmail.Thread
	pageToken := ""

	// 1. Fetch all threads with srclabel.JobsOpportunity
	// 2. Filter all messages before a reply
	// 3. Pick the earliest message
	// 4. Hit the parse endpoint
	// 5. if all fields are present, insert a job
	// 6. if any fields are missing, log

	// batch process messages to reduce memory usage
	for {
		// get next set of messages
		// if this is the first notification, trigger a one-time sync
		threads, pageToken, err = fetchJobThreadsSinceDate(srv, startDate, pageToken)

		// for now, abort on error
		if err != nil {
			handleError(w, "error fetching emails", err)
			return
		}

		// get the messages for each thread
		messages := map[string]*gmail.Message{}

		for _, t := range threads {
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

		// predict one at a time for now
		for id, input := range inputs {
			log.Printf("parsing email: %s", id)
			job, err := mlSrv.ParseJob(input)
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
				UserID:        user.UserID,
				UserEmail:     email,
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

	log.Println("done.")
}

// fetchJobThreadsSinceDate fetches all job threads since the start date
func fetchJobThreadsSinceDate(srv *srcmail.Service, date time.Time, pageToken string) ([]*gmail.Thread, string, error) {
	q := fmt.Sprintf("-label:sent label:%s after:%s", srclabel.JobsOpportunity.Name, date.Format("2006/01/02"))

	r, err := srcmail.ExecuteWithRetries(func() (*gmail.ListThreadsResponse, error) {
		return srv.Users.Threads.
			List(srv.UserID).
			PageToken(pageToken).
			Q(q).
			MaxResults(100).
			Do()
	})

	if err != nil {
		return nil, "", err
	}

	return r.Threads, r.NextPageToken, nil
}
