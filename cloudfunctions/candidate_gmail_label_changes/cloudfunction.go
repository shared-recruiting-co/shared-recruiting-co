package cloudfunctions

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/idtoken"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
	srcmail "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail"
	srclabel "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/label"
	srcmessage "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/mail/gmail/message"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

const (
	provider = "google"
)

func init() {
	functions.CloudEvent("Handler", handler)
}

// naive error handler for now
func handleError(msg string, err error) error {
	err = fmt.Errorf("%s: %w", msg, err)
	sentry.CaptureException(err)
	return err
}

func jsonFromEnv(env string) ([]byte, error) {
	encoded := os.Getenv(env)
	decoded, err := base64.URLEncoding.DecodeString(encoded)

	return decoded, err
}

func contains[T comparable](list []T, item T) bool {
	for _, element := range list {
		if element == item {
			return true
		}
	}
	return false
}

type CloudFunction struct {
	ctx                  context.Context
	queries              db.Querier
	srv                  *srcmail.Service
	labels               *srclabel.CandidateLabels
	model                ml.Service
	user                 db.UserProfile
	examplesCollectorSrv *srcmail.Service
	payload              *schema.EmailLabelChanges
	addedLabelFuncs      map[string]func(cf *CloudFunction, msg *gmail.Message) error
	removedLabelFuncs    map[string]func(cf *CloudFunction, msg *gmail.Message) error
}

func NewCloudFunction(ctx context.Context, payload schema.EmailLabelChanges) (*CloudFunction, error) {
	creds, err := jsonFromEnv("GOOGLE_OAUTH2_CREDENTIALS")
	if err != nil {
		return nil, fmt.Errorf("error parsing GOOGLE_OAUTH2_CREDENTIALScredentials: %w", err)
	}

	// 0, Create SRC client
	apiURL := os.Getenv("SUPABASE_API_URL")
	apiKey := os.Getenv("SUPABASE_API_KEY")
	queries := db.NewHTTP(apiURL, apiKey)

	// 1. Get User from email address
	user, err := queries.GetUserProfileByEmail(ctx, payload.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile by email: %w", err)
	}

	var examplesCollectorSrv *srcmail.Service
	if user.AutoContribute {
		auth, err := jsonFromEnv("EXAMPLES_GMAIL_OAUTH_TOKEN")
		if err != nil {
			return nil, fmt.Errorf("error parsing EXAMPLES_GMAIL_OAUTH_TOKEN: %w", err)
		}
		examplesCollectorSrv, err = srcmail.NewService(ctx, creds, auth)
		if err != nil {
			return nil, fmt.Errorf("error creating examples collector service: %w", err)
		}
	}

	// 2. Get User' OAuth Token
	userToken, err := queries.GetUserOAuthToken(ctx, db.GetUserOAuthTokenParams{
		UserID:   user.UserID,
		Email:    payload.Email,
		Provider: provider,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting user oauth token: %w", err)
	}

	// stop early if user token is marked invalid
	if !userToken.IsValid {
		return nil, fmt.Errorf("user token is not valid: %s", payload.Email)
	}

	// 3. Create Gmail Service
	auth := []byte(userToken.Token)
	srv, err := srcmail.NewService(ctx, creds, auth)
	if err != nil {
		return nil, fmt.Errorf("error creating gmail service: %w", err)
	}

	// 4. Get or Create SRC Labels
	labels, err := srv.GetOrCreateCandidateLabels()
	if err != nil {
		// first request, so check if the error is an oauth error
		// if so, update the database
		if srcmail.IsOAuth2Error(err) {
			log.Printf("error oauth error: %v", err)
			// update the user's oauth token
			err = queries.UpsertUserOAuthToken(ctx, db.UpsertUserOAuthTokenParams{
				UserID:   userToken.UserID,
				Email:    payload.Email,
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
		return nil, fmt.Errorf("error getting or creating src labels: %w", err)
	}
	// Create recruiting detector client
	mlServiceBaseURL := os.Getenv("ML_SERVICE_URL")
	idTokenSource, err := idtoken.NewTokenSource(ctx, mlServiceBaseURL)
	if err != nil {
		return nil, fmt.Errorf("error creating id token source: %w", err)
	}

	idToken, err := idTokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("error getting id token: %w", err)
	}

	model := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   mlServiceBaseURL,
		AuthToken: idToken.AccessToken,
	})

	// set up label handlers
	// use label IDs as keys
	addedLabelFuncs := map[string]func(cf *CloudFunction, msg *gmail.Message) error{
		labels.JobsOpportunity.Id:   handleAddedJobOpportunityLabel,
		labels.JobsInterested.Id:    handleAddedJobInterestedLabel,
		labels.JobsNotInterested.Id: handleAddedJobNotInterestedLabel,
		labels.JobsSaved.Id:         handleAddedJobSavedLabel,
	}
	removedLabelFuncs := map[string]func(cf *CloudFunction, msg *gmail.Message) error{
		labels.JobsOpportunity.Id:   handleRemovedJobOpportunityLabel,
		labels.JobsInterested.Id:    handleRemovedJobInterestedLabel,
		labels.JobsNotInterested.Id: handleRemovedJobNotInterestedLabel,
		labels.JobsSaved.Id:         handleRemovedJobSavedLabel,
	}

	return &CloudFunction{
		ctx:                  ctx,
		queries:              queries,
		srv:                  srv,
		labels:               labels,
		model:                model,
		user:                 user,
		examplesCollectorSrv: examplesCollectorSrv,
		payload:              &payload,
		addedLabelFuncs:      addedLabelFuncs,
		removedLabelFuncs:    removedLabelFuncs,
	}, nil
}

func handler(ctx context.Context, e event.Event) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       os.Getenv("FUNCTION_NAME"),
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %v", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	var msg schema.MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return handleError("error parsing Pub/Sub message", err)
	}

	data := msg.Message.Data
	log.Printf("Event: %s", data)

	var payload schema.EmailLabelChanges
	if err := json.Unmarshal(data, &payload); err != nil {
		return handleError("error parsing email messages", err)
	}

	// validate payload
	// for invalid payloads, we don't want to retry
	if payload.Email == "" || len(payload.Changes) == 0 {
		err = fmt.Errorf("received invalid payload: %v", payload)
		log.Print(err)
		sentry.CaptureException(err)
		return nil
	}

	cf, err := NewCloudFunction(ctx, payload)
	if err != nil {
		return handleError("error creating cloud function", err)
	}

	err = cf.processLabelChanges(payload.Changes)
	if err != nil {
		return handleError("error processing messages", err)
	}

	return nil
}

func (cf *CloudFunction) processLabelChanges(changes []schema.EmailLabelChange) error {
	// placeholder logic
	// for now, we'll just log the changes
	for _, change := range changes {
		log.Printf("change: %v", change)
		// 1. get the message
		msg, err := cf.srv.GetMessage(change.MessageID)
		if err != nil {
			log.Printf("error getting message: %v", err)
			continue
		}
		if change.ChangeType == schema.EmailLabelChangeTypeAdded {
			// 2. for each label, check if we have a function for it
			for _, label := range change.LabelIDs {
				// Check the label is still added
				if !contains(msg.LabelIds, label) {
					log.Printf("label %s was removed from message %s before processing", label, change.MessageID)
					continue
				}
				// process
				err = cf.processLabelAdded(msg, label)
				// abort on error so we can retry
				if err != nil {
					return fmt.Errorf("error processing label added: %w", err)
				}
			}
		} else if change.ChangeType == schema.EmailLabelChangeTypeRemoved {
			// Check the label is still removed
			for _, label := range change.LabelIDs {
				if contains(msg.LabelIds, label) {
					log.Printf("label %s was added back to message %s before processing", label, change.MessageID)
					continue
				}
				// process
				err = cf.processLabelRemoved(msg, label)
				// abort on error so we can retry
				if err != nil {
					return fmt.Errorf("error processing label removed: %w", err)
				}
			}
		} else {
			log.Printf("unknown change type: %v", change.ChangeType)
		}
	}

	return nil
}

func (cf *CloudFunction) processLabelAdded(msg *gmail.Message, labelID string) error {
	// 1. check if we have a function for the label
	fn, ok := cf.addedLabelFuncs[labelID]
	if !ok {
		log.Printf("no add function found for label %s", labelID)
		return nil
	}
	// 2. run the function
	return fn(cf, msg)
}

func (cf *CloudFunction) processLabelRemoved(msg *gmail.Message, labelID string) error {
	// 1. check if we have a function for the label
	fn, ok := cf.removedLabelFuncs[labelID]
	if !ok {
		log.Printf("no remove function found for label %s", labelID)
		return nil
	}
	// 2. run the function
	return fn(cf, msg)
}

func (cf *CloudFunction) ParseEmail(msg *gmail.Message) (*ml.ParseJobResponse, error) {
	parseRequest := ml.ParseJobRequest{
		From:    srcmessage.Sender(msg),
		Subject: srcmessage.Subject(msg),
		Body:    srcmessage.Body(msg),
	}
	log.Printf("parsing email: %s", msg.Id)
	return cf.model.ParseJob(&parseRequest)
}

func (cf *CloudFunction) InsertRecruiterEmailIntoDB(msg *gmail.Message, company, title, recruiter string) error {
	recruiterEmail := srcmessage.SenderEmail(msg)
	data := map[string]interface{}{
		"recruiter":       recruiter,
		"recruiter_email": recruiterEmail,
	}

	// turn data into json.RawMessage
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return cf.queries.InsertUserEmailJob(cf.ctx, db.InsertUserEmailJobParams{
		UserID:        cf.user.UserID,
		UserEmail:     cf.payload.Email,
		EmailThreadID: msg.ThreadId,
		// convert epoch ms to time.Time
		EmailedAt: srcmessage.CreatedAt(msg),
		Company:   company,
		JobTitle:  title,
		Data:      b,
	})
}

func (cf *CloudFunction) IsKnownRecruitingEmail(msg *gmail.Message) bool {
	// check if message is a known recruiting outbound message
	internalMessageID := srcmessage.Header(msg, "Message-ID")
	if internalMessageID == "" {
		log.Printf("skipping message: %s, no Message-ID header", msg.Id)
		return false
	}
	// check if message is a known recruiting outbound message
	// look up by RFC2822 message ID and to_email
	// TODO: Is it better to use the recipient email from the header or the current user?
	recipient := srcmessage.RecipientEmail(msg)
	if recipient == "" {
		log.Printf("skipping message: %s, no recipient email", msg.Id)
		return false
	}

	_, err := cf.queries.GetRecruiterOutboundMessageByRecipient(cf.ctx, db.GetRecruiterOutboundMessageByRecipientParams{
		ToEmail:           recipient,
		InternalMessageID: internalMessageID,
	})
	// we expect a not found error
	if err == nil {
		// log and add to known messages
		log.Printf("found known recruiting message: %s", msg.Id)
		return true
	} else if err != sql.ErrNoRows {
		log.Printf("error looking up known recruiting message: %v", err)
	}
	return false
}

// GetJobIDForMessage gets the job ID for a gmail message
// TODO
// Right now we make multiple requests to handle both user_email_job and candidate_company_inbound
// We need to consolidate these into a single query
func (cf *CloudFunction) GetJobIDForMessage(msg *gmail.Message) (*uuid.UUID, error) {
	// - first check user_email_job
	emailJob, err := cf.queries.GetUserEmailJobByThreadID(cf.ctx, db.GetUserEmailJobByThreadIDParams{
		UserEmail:     cf.payload.Email,
		EmailThreadID: msg.ThreadId,
	})

	// success
	if err == nil {
		return &emailJob.JobID, nil
	} else if err != sql.ErrNoRows {
		// unexpected error
		return nil, err
	}

	// check recruiter_outbound_message
	internalMessageID := srcmessage.Header(msg, "Message-ID")
	if internalMessageID == "" {
		log.Printf("skipping message: %s, no Message-ID header", msg.Id)
		return nil, nil
	}
	// check if message is a known recruiting outbound message
	// look up by RFC2822 message ID and to_email
	// TODO: Is it better to use the recipient email from the header or the current user?
	recipient := srcmessage.RecipientEmail(msg)
	if recipient == "" {
		log.Printf("skipping message: %s, no recipient email", msg.Id)
		return nil, nil
	}

	outbound, err := cf.queries.GetRecruiterOutboundMessageByRecipient(cf.ctx, db.GetRecruiterOutboundMessageByRecipientParams{
		ToEmail:           recipient,
		InternalMessageID: internalMessageID,
	})
	if err == sql.ErrNoRows {
		// not found
		return nil, nil
	} else if err != nil {
		// unexpected error
		return nil, err
	} else if !outbound.TemplateID.Valid {
		// no template ID
		return nil, nil
	}

	// use template ID to get job ID
	template, err := cf.queries.GetRecruiterOutboundTemplate(cf.ctx, outbound.TemplateID.UUID)
	if err != nil {
		return nil, err
	}
	// check if job ID is valid
	if !template.JobID.Valid {
		log.Printf("template %s has no job ID", outbound.TemplateID.UUID.String())
		return nil, nil
	}

	return &template.JobID.UUID, nil
}

// Label Functions

// handleAddedJobOpportunityLabel handles the added job opportunities label
// 1. check if the user has auto contribute enabled
// 2. parse the email and add to the user's job board
func handleAddedJobOpportunityLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("added job opportunities label: %s", msg.Id)

	// 1. check if the user has auto contribute enabled
	if cf.user.AutoContribute && cf.examplesCollectorSrv != nil {
		// clone the message to the examples inbox
		_, err := srcmail.CloneMessage(cf.srv, cf.examplesCollectorSrv, msg.Id, []string{"INBOX", "UNREAD"})

		if err != nil {
			// don't abort on error
			log.Printf("error collecting email %s: %v", msg.Id, err)
			sentry.CaptureException(fmt.Errorf("error collecting email %s: %w", msg.Id, err))
		}
	}

	// 2. parse the email and add to the user's job board
	// only parse jobs if the msg is NOT a verified job
	if !cf.IsKnownRecruitingEmail(msg) {
		job, err := cf.ParseEmail(msg)
		// for now, abort on error
		if err != nil {
			return err
		}

		// if fields are missing, skip
		if job.Company == "" || job.Title == "" || job.Recruiter == "" {
			// print sender and subject
			log.Printf("skipping job: %v", job)
			return nil
		}

		err = cf.InsertRecruiterEmailIntoDB(msg, job.Company, job.Title, job.Recruiter)

		// for now, continue on error
		if err != nil {
			log.Printf("error inserting job (%v): %v", job, err)
		}

	}

	return nil
}

// handleRemovedJobOpportunityLabel handles the removed job opportunities label
// 1. remove from a user's user_email_job if it exists
func handleRemovedJobOpportunityLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("removed job opportunities label: %s", msg.Id)

	// Remove from a user's user_email_job if it exists
	err := cf.queries.DeleteUserEmailJobByEmailThreadID(cf.ctx, db.DeleteUserEmailJobByEmailThreadIDParams{
		UserEmail:     cf.payload.Email,
		EmailThreadID: msg.ThreadId,
	})
	if err != nil {
		log.Printf("error deleting job: %v", err)
	}

	return nil
}

// handleAddedJobInterestedLabel handles the added job interest label
func handleAddedJobInterestedLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("added job interested label: %s", msg.Id)

	// update message to remove other job interest labels
	_, err := srcmail.ExecuteWithRetries(func() (interface{}, error) {
		return cf.srv.Users.Messages.Modify(cf.srv.UserID, msg.Id, &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{cf.labels.JobsNotInterested.Id, cf.labels.JobsSaved.Id},
		}).Do()
	})
	if err != nil {
		// for now, log and continue
		log.Printf("error removing other job interest labels (interested): %v", err)
		sentry.CaptureException(fmt.Errorf("error removing other job interest labels (interested): %w", err))
	}
	// update database
	// 1. Get the job from the database
	jobID, err := cf.GetJobIDForMessage(msg)
	if err != nil {
		log.Printf("error getting job id for message: %v", err)
		return nil
	} else if jobID == nil {
		log.Printf("no job id for message: %v", msg.Id)
		return nil
	}

	// 2. Update the job with the interested label
	err = cf.queries.UpsertCandidateJobInterest(cf.ctx, db.UpsertCandidateJobInterestParams{
		CandidateID: cf.user.UserID,
		JobID:       *jobID,
		Interest:    db.JobInterestInterested,
	})

	// for now, only log errors
	if err != nil {
		log.Printf("error updating job interest: %v", err)
		sentry.CaptureException(fmt.Errorf("error updating job interest: %w", err))
	}

	return nil
}

// handleRemovedJobInterestedLabel handles the removed job interest label
func handleRemovedJobInterestedLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("removed job interested label: %s", msg.Id)

	// update database (set to null if set to 'interest')
	jobID, err := cf.GetJobIDForMessage(msg)
	if err != nil {
		log.Printf("error getting job id for message: %v", err)
		return nil
	} else if jobID == nil {
		log.Printf("no job id for message: %v", msg.Id)
		return nil
	}

	err = cf.queries.DeleteCandidateJobInterestConditionally(cf.ctx, db.DeleteCandidateJobInterestConditionallyParams{
		CandidateID: cf.user.UserID,
		JobID:       *jobID,
		Interest:    db.JobInterestInterested,
	})

	// for now, only log errors
	if err != nil {
		log.Printf("error updating job interest: %v", err)
		sentry.CaptureException(fmt.Errorf("error updating job interest: %w", err))
	}

	return nil
}

// handleAddedJobNotInterestedLabel handles the job not interested label
func handleAddedJobNotInterestedLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("added job not interested label: %s", msg.Id)

	// update message to remove other job interest labels
	_, err := srcmail.ExecuteWithRetries(func() (interface{}, error) {
		return cf.srv.Users.Messages.Modify(cf.srv.UserID, msg.Id, &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{cf.labels.JobsInterested.Id, cf.labels.JobsSaved.Id},
		}).Do()
	})
	if err != nil {
		// for now, log and continue
		log.Printf("error removing other job interest labels (not interested): %v", err)
		sentry.CaptureException(fmt.Errorf("error removing other job interest labels (not interested): %w", err))
	}
	// update database
	// 1. Get the job from the database
	jobID, err := cf.GetJobIDForMessage(msg)
	if err != nil {
		log.Printf("error getting job id for message: %v", err)
		return nil
	} else if jobID == nil {
		log.Printf("no job id for message: %v", msg.Id)
		return nil
	}

	// 2. Update the job with the interested label
	err = cf.queries.UpsertCandidateJobInterest(cf.ctx, db.UpsertCandidateJobInterestParams{
		CandidateID: cf.user.UserID,
		JobID:       *jobID,
		Interest:    db.JobInterestNotInterested,
	})

	// for now, only log errors
	if err != nil {
		log.Printf("error updating job interest: %v", err)
		sentry.CaptureException(fmt.Errorf("error updating job interest: %w", err))
	}

	return nil
}

// handleRemovedJobNotInterestedLabel handles the removed job not interested label
func handleRemovedJobNotInterestedLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("removed job not interested label: %s", msg.Id)

	// update database (set to null if set to 'not_interest')
	jobID, err := cf.GetJobIDForMessage(msg)
	if err != nil {
		log.Printf("error getting job id for message: %v", err)
		return nil
	} else if jobID == nil {
		log.Printf("no job id for message: %v", msg.Id)
		return nil
	}

	err = cf.queries.DeleteCandidateJobInterestConditionally(cf.ctx, db.DeleteCandidateJobInterestConditionallyParams{
		CandidateID: cf.user.UserID,
		JobID:       *jobID,
		Interest:    db.JobInterestNotInterested,
	})

	// for now, only log errors
	if err != nil {
		log.Printf("error updating job interest: %v", err)
		sentry.CaptureException(fmt.Errorf("error updating job interest: %w", err))
	}

	return nil
}

// handleAddedJobSavedLabel handles the job saved label
func handleAddedJobSavedLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("added job saved label: %s", msg.Id)

	// update message to remove other job interest labels
	_, err := srcmail.ExecuteWithRetries(func() (interface{}, error) {
		return cf.srv.Users.Messages.Modify(cf.srv.UserID, msg.Id, &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{cf.labels.JobsInterested.Id, cf.labels.JobsNotInterested.Id},
		}).Do()
	})
	if err != nil {
		// for now, log and continue
		log.Printf("error removing other job interest labels (saved): %v", err)
		sentry.CaptureException(fmt.Errorf("error removing job interest labels (saved): %w", err))
	}

	// update database
	// 1. Get the job from the database
	jobID, err := cf.GetJobIDForMessage(msg)
	if err != nil {
		log.Printf("error getting job id for message: %v", err)
		return nil
	} else if jobID == nil {
		log.Printf("no job id for message: %v", msg.Id)
		return nil
	}

	// 2. Update the job with the interested label
	err = cf.queries.UpsertCandidateJobInterest(cf.ctx, db.UpsertCandidateJobInterestParams{
		CandidateID: cf.user.UserID,
		JobID:       *jobID,
		Interest:    db.JobInterestSaved,
	})

	// for now, only log errors
	if err != nil {
		log.Printf("error updating job interest: %v", err)
		sentry.CaptureException(fmt.Errorf("error updating job interest: %w", err))
	}

	return nil
}

// handleRemovedJobSavedLabel handles the removed job saved label
func handleRemovedJobSavedLabel(cf *CloudFunction, msg *gmail.Message) error {
	// For now log
	log.Printf("removed job saved label: %s", msg.Id)

	// update database (set to null if set to 'saved')
	jobID, err := cf.GetJobIDForMessage(msg)
	if err != nil {
		log.Printf("error getting job id for message: %v", err)
		return nil
	} else if jobID == nil {
		log.Printf("no job id for message: %v", msg.Id)
		return nil
	}

	err = cf.queries.DeleteCandidateJobInterestConditionally(cf.ctx, db.DeleteCandidateJobInterestConditionallyParams{
		CandidateID: cf.user.UserID,
		JobID:       *jobID,
		Interest:    db.JobInterestSaved,
	})

	// for now, only log errors
	if err != nil {
		log.Printf("error updating job interest: %v", err)
		sentry.CaptureException(fmt.Errorf("error updating job interest: %w", err))
	}

	return nil
}
