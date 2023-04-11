package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// HTTPQueries is a client for the PostgREST API
type HTTPQueries struct {
	client *http.Client
	// URL is the Base URL of the PostgREST API
	URL string
	// APIKey is the API Key for the PostgREST API
	APIKey string
	// Debug enables debug logging
	Debug bool
}

// NewHTTP creates a new HTTPQueries.
func NewHTTP(url, apiKey string) *HTTPQueries {
	return &HTTPQueries{
		client: &http.Client{},
		URL:    url,
		APIKey: apiKey,
	}
}

// sanitize user input according to: ://github.com/shared-recruiting-co/shared-recruiting-co/security/code-scanning/8
func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// DoRequest performs a request to the PostgREST API.
func (q *HTTPQueries) DoRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	reqPath, err := url.JoinPath(q.URL, path)
	if err != nil {
		return nil, fmt.Errorf("error joining path: %w", err)
	}
	// url.JoinPath escapes the query string, so we need to unescape it
	reqPath, err = url.QueryUnescape(reqPath)
	if err != nil {
		return nil, fmt.Errorf("error unescaping path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqPath, body)
	if err != nil {
		return nil, err
	}
	if q.Debug {
		log.Printf("request: %s %s\n", req.Method, sanitize(req.URL.String()))
	}

	req.Header.Set("apikey", q.APIKey)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", q.APIKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Range-Unit", "items")
	// always upsert on POST (aka insert)
	if method == http.MethodPost {
		req.Header.Add("Prefer", "resolution=merge-duplicates")
		// always return the inserted object
		req.Header.Add("Prefer", "return=representation")
	}

	return q.client.Do(req)
}

// singleOrError is a helper function to return a single entry or an error.
// By default PostgREST always returns an array even if there is only one entry.
func singleOrError[T any](slice []T) (T, error) {
	var result T
	if len(slice) == 0 {
		// for now, return same error as database client
		return result, sql.ErrNoRows
	}

	if len(slice) > 1 {
		return result, fmt.Errorf("more than one element in slice")
	}

	return slice[0], nil
}

// GetUserProfileByEmail fetches a user profile by email.
func (q *HTTPQueries) GetUserProfileByEmail(ctx context.Context, email string) (UserProfile, error) {
	basePath := "/rpc/get_user_profile_by_email"
	path := basePath
	var result UserProfile

	body, err := json.Marshal(map[string]string{"input": email})
	if err != nil {
		return result, err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("fetch user profile by email: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	if result.Email == "" {
		return result, sql.ErrNoRows
	}

	return result, err
}

// GetUserEmailSyncHistory fetches a user's email sync history.
func (q *HTTPQueries) GetUserEmailSyncHistory(ctx context.Context, arg GetUserEmailSyncHistoryParams) (UserEmailSyncHistory, error) {
	basePath := "/user_email_sync_history"
	query := fmt.Sprintf("select=*&user_id=eq.%s&inbox_type=eq.%s&email=eq.%s", arg.UserID, arg.InboxType, arg.Email)
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result UserEmailSyncHistory

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("fetch user email sync history: %s", resp.Status)
	}

	var history []UserEmailSyncHistory
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return result, err
	}

	result, err = singleOrError(history)
	return result, err
}

// GetUserOAuthToken fetches a user's oauth token.
func (q *HTTPQueries) GetUserOAuthToken(ctx context.Context, arg GetUserOAuthTokenParams) (UserOauthToken, error) {
	basePath := "/user_oauth_token"
	query := fmt.Sprintf("select=*&user_id=eq.%s&email=eq.%s&provider=eq.%s", arg.UserID, arg.Email, arg.Provider)
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result UserOauthToken

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("fetch user oauth token: %s", resp.Status)
	}

	var tokens []UserOauthToken
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return result, err
	}

	result, err = singleOrError(tokens)
	return result, err
}

// UpsertUserOAuthToken
func (q *HTTPQueries) UpsertUserOAuthToken(ctx context.Context, arg UpsertUserOAuthTokenParams) error {
	basePath := "/user_oauth_token"
	path := basePath
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upsert user email sync history: %s", resp.Status)
	}

	return nil
}

// ListUserOAuthTokens lists a user oauth tokens.
func (q *HTTPQueries) ListUserOAuthTokens(ctx context.Context, arg ListUserOAuthTokensParams) ([]UserOauthToken, error) {
	basePath := "/user_oauth_token"
	query := "select=*"
	if arg.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	query = fmt.Sprintf("%s&provider=eq.%s", query, arg.Provider)
	query = fmt.Sprintf("%s&is_valid=eq.%t", query, arg.IsValid)
	// pagination params
	query = fmt.Sprintf("%s&limit=%d&offset=%d", query, arg.Limit, arg.Offset)
	query = fmt.Sprintf("%s&limit=%d&offset=%d", query, arg.Limit, arg.Offset)

	path := fmt.Sprintf("%s?%s", basePath, query)
	var result []UserOauthToken

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("list valid user oauth tokens: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	if len(result) == 0 {
		// for now, return same error as database client
		return result, sql.ErrNoRows
	}

	return result, nil
}

// ListCandidateOAuthTokens lists a user oauth tokens.
func (q *HTTPQueries) ListCandidateOAuthTokens(ctx context.Context, arg ListCandidateOAuthTokensParams) ([]CandidateOauthToken, error) {
	basePath := "/candidate_oauth_token"
	query := "select=*"
	if arg.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	query = fmt.Sprintf("%s&provider=eq.%s", query, arg.Provider)
	query = fmt.Sprintf("%s&is_valid=eq.%t", query, arg.IsValid)
	// pagination params
	query = fmt.Sprintf("%s&limit=%d&offset=%d", query, arg.Limit, arg.Offset)

	path := fmt.Sprintf("%s?%s", basePath, query)
	var result []CandidateOauthToken

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("list valid candidate oauth tokens: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	if len(result) == 0 {
		// for now, return same error as database client
		return result, sql.ErrNoRows
	}

	return result, nil
}

// ListRecruiterOAuthTokens lists a user oauth tokens.
func (q *HTTPQueries) ListRecruiterOAuthTokens(ctx context.Context, arg ListRecruiterOAuthTokensParams) ([]RecruiterOauthToken, error) {
	basePath := "/recruiter_oauth_token"
	query := "select=*"
	if arg.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	query = fmt.Sprintf("%s&provider=eq.%s", query, arg.Provider)
	query = fmt.Sprintf("%s&is_valid=eq.%t", query, arg.IsValid)
	// pagination params
	query = fmt.Sprintf("%s&limit=%d&offset=%d", query, arg.Limit, arg.Offset)

	path := fmt.Sprintf("%s?%s", basePath, query)
	var result []RecruiterOauthToken

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("list valid user oauth tokens: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	if len(result) == 0 {
		// for now, return same error as database client
		return result, sql.ErrNoRows
	}

	return result, nil
}

// UpsertUserEmailSyncHistory upserts a user's email sync history.
func (q *HTTPQueries) UpsertUserEmailSyncHistory(ctx context.Context, arg UpsertUserEmailSyncHistoryParams) error {
	basePath := "/user_email_sync_history"
	path := basePath
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upsert user email sync history: %s", resp.Status)
	}

	return nil
}

// IncrementUserEmailStat increments a user's email stat.
func (q *HTTPQueries) IncrementUserEmailStat(ctx context.Context, arg IncrementUserEmailStatParams) error {
	// we cannot do upserts with PostgREST, so instead we user a stored procedure
	basePath := "/rpc/increment_user_email_stat"
	path := basePath
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// RPCs return 204
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("increment user email stat: %s", resp.Status)
	}

	return nil
}

// GetUserEmailJob fetches a user's email job by job ID
func (q *HTTPQueries) GetUserEmailJob(ctx context.Context, jobID uuid.UUID) (UserEmailJob, error) {
	basePath := "/user_email_job"
	query := fmt.Sprintf("select=*&job_id=eq.%s", jobID.String())
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result UserEmailJob

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error fetching user email job: %s", resp.Status)
	}

	var data []UserEmailJob
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return result, err
	}

	result, err = singleOrError(data)
	return result, err
}

// GetUserEmailJobByThreadID fetches a user's email job by user email and thread ID
func (q *HTTPQueries) GetUserEmailJobByThreadID(ctx context.Context, arg GetUserEmailJobByThreadIDParams) (UserEmailJob, error) {
	basePath := "/user_email_job"
	query := fmt.Sprintf("select=*&user_email=eq.%s&email_thread_id=eq.%s", arg.UserEmail, arg.EmailThreadID)
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result UserEmailJob

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error fetching user email job by thread ID: %s", resp.Status)
	}

	var data []UserEmailJob
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return result, err
	}

	result, err = singleOrError(data)
	return result, err
}

// ListUserEmailJobs lists a user's email jobs.
func (q *HTTPQueries) ListUserEmailJobs(ctx context.Context, arg ListUserEmailJobsParams) ([]UserEmailJob, error) {
	basePath := "/user_email_job"
	query := "select=*&order=emailed_at.desc"
	query = fmt.Sprintf("%s&user_id=eq.%s&limit=%d&offset=%d", query, arg.UserID.String(), arg.Limit, arg.Offset)

	path := fmt.Sprintf("%s?%s", basePath, query)
	var result []UserEmailJob

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error listing user email jobs: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	if len(result) == 0 {
		// for now, return same error as database client
		return result, sql.ErrNoRows
	}

	return result, nil
}

// InsertUserEmailJob inserts a user's email job.
func (q *HTTPQueries) InsertUserEmailJob(ctx context.Context, arg InsertUserEmailJobParams) error {
	basePath := "/user_email_job"
	path := basePath
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error inserting user email job: %s", resp.Status)
	}

	return nil
}

// DeleteUserEmailJobByEmailThreadID deletes a user's email job by email thread ID.
func (q *HTTPQueries) DeleteUserEmailJobByEmailThreadID(ctx context.Context, arg DeleteUserEmailJobByEmailThreadIDParams) error {
	basePath := "/user_email_job"
	query := fmt.Sprintf("user_email=eq.%s&email_thread_id=eq.%s", arg.UserEmail, arg.EmailThreadID)
	path := fmt.Sprintf("%s?%s", basePath, query)

	resp, err := q.DoRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting user email job: %s", resp.Status)
	}

	return nil
}

// CountUserEmailJobs counts the number of user's email jobs.
func (q *HTTPQueries) CountUserEmailJobs(ctx context.Context, userID uuid.UUID) (int64, error) {
	basePath := "/user_email_job"
	query := fmt.Sprintf("user_id=eq.%s", userID.String())
	path := fmt.Sprintf("%s?%s", basePath, query)

	resp, err := q.DoRequest(ctx, http.MethodHead, path, nil)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusPartialContent {
		return 0, fmt.Errorf("error counting user email jobs: %s", resp.Status)
	}

	// parse count from header
	contentedRange := resp.Header.Get("Content-Range")
	parts := strings.Split(contentedRange, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("error parsing count from Content-Range header: %s", contentedRange)
	}
	count, err := strconv.ParseInt(parts[1], 10, 64)
	return count, err
}

// GetRecruiterByEmail fetches a recruiter profile given their email
func (q *HTTPQueries) GetRecruiterByEmail(ctx context.Context, email string) (GetRecruiterByEmailRow, error) {
	basePath := "/rpc/get_recruiter_by_email"
	path := basePath
	var result GetRecruiterByEmailRow

	body, err := json.Marshal(map[string]string{"input": email})
	if err != nil {
		return result, err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error fetching recruiter by email: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	if result.Email == "" {
		return result, sql.ErrNoRows
	}

	return result, err
}

// GetRecruiterOutboundMessage fetches a recruiter's outbound message by message ID
func (q *HTTPQueries) GetRecruiterOutboundMessage(ctx context.Context, arg GetRecruiterOutboundMessageParams) (RecruiterOutboundMessage, error) {
	basePath := "/recruiter_outbound_message"
	query := fmt.Sprintf("select=*&recruiter_id=eq.%s&message_id=eq.%s", arg.RecruiterID.String(), arg.MessageID)
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result RecruiterOutboundMessage

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error fetching recruiter outbound message: %s", resp.Status)
	}

	var results []RecruiterOutboundMessage
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return result, err
	}

	result, err = singleOrError(results)
	return result, err
}

// GetRecruiterOutboundMessageByRecipient fetches a recruiter's outbound message by message ID
func (q *HTTPQueries) GetRecruiterOutboundMessageByRecipient(ctx context.Context, arg GetRecruiterOutboundMessageByRecipientParams) (RecruiterOutboundMessage, error) {
	basePath := "/recruiter_outbound_message"
	query := fmt.Sprintf("select=*&to_email=eq.%s&internal_message_id=eq.%s", arg.ToEmail, arg.InternalMessageID)
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result RecruiterOutboundMessage

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error fetching recruiter outbound message by recipient: %s", resp.Status)
	}

	var results []RecruiterOutboundMessage
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return result, err
	}

	result, err = singleOrError(results)
	return result, err
}

// InsertRecruiterOutboundMessage inserts a recruiter's outbound message
func (q *HTTPQueries) InsertRecruiterOutboundMessage(ctx context.Context, arg InsertRecruiterOutboundMessageParams) error {
	basePath := "/recruiter_outbound_message"
	path := basePath
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("insert recruiter outbound message: %s", resp.Status)
	}

	return nil
}

func (q *HTTPQueries) GetRecruiterOutboundTemplate(ctx context.Context, templateID uuid.UUID) (RecruiterOutboundTemplate, error) {
	basePath := "/recruiter_outbound_template"
	query := fmt.Sprintf("select=*&template_id=eq.%s", templateID.String())
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result RecruiterOutboundTemplate

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("error fetching recruiter outbound template: %s", resp.Status)
	}

	var results []RecruiterOutboundTemplate
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return result, err
	}

	result, err = singleOrError(results)
	return result, err
}

// InsertRecruiterOutboundTemplate inserts a recruiter's outbound template
func (q *HTTPQueries) InsertRecruiterOutboundTemplate(ctx context.Context, arg InsertRecruiterOutboundTemplateParams) (RecruiterOutboundTemplate, error) {
	basePath := "/recruiter_outbound_template"
	path := basePath
	body, err := json.Marshal(arg)
	var result RecruiterOutboundTemplate
	if err != nil {
		return result, fmt.Errorf("error marshalling recruiter outbound template: %w\n%v", err, arg)
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return result, fmt.Errorf("insert recruiter outbound template: %s", resp.Status)
	}

	// POST returns an array of the inserted rows
	var results []RecruiterOutboundTemplate
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return result, fmt.Errorf("error decoding recruiter outbound template: %w", err)
	}
	if len(results) != 1 {
		return result, fmt.Errorf("unexpected number of recruiter outbound templates: %d\n%v", len(results), results)
	}
	result = results[0]
	return result, nil
}

func (q *HTTPQueries) ListSimilarRecruiterOutboundTemplates(ctx context.Context, arg ListSimilarRecruiterOutboundTemplatesParams) ([]ListSimilarRecruiterOutboundTemplatesRow, error) {
	basePath := "/rpc/list_similar_recruiter_outbound_templates"
	path := basePath
	var results []ListSimilarRecruiterOutboundTemplatesRow

	body, err := json.Marshal(arg)
	if err != nil {
		return results, err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return results, fmt.Errorf("fetch user profile by email: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return results, err
	}

	// no results is not an error here
	return results, err
}

// UpsertCandidateJobInterest upserts a candidate's job interest
func (q *HTTPQueries) UpsertCandidateJobInterest(ctx context.Context, arg UpsertCandidateJobInterestParams) error {
	basePath := "/candidate_job_interest"
	path := basePath
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error upserting candidate job interest: %s", resp.Status)
	}

	return nil
}

type UpdateCandidateJobInterestBody struct {
	Interest NullJobInterest `json:"interest"`
}

func (q *HTTPQueries) UpdateCandidateJobInterestConditionally(ctx context.Context, arg UpdateCandidateJobInterestConditionallyParams) error {
	basePath := "/candidate_job_interest"
	query := fmt.Sprintf("candidate_id=eq.%s&job_id=eq.%s&interest=eq.%s", arg.CandidateID.String(), arg.JobID.String(), arg.Interest)
	path := fmt.Sprintf("%s?%s", basePath, query)
	body, err := json.Marshal(UpdateCandidateJobInterestBody{Interest: arg.SetInterest})
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error conditionally updating candidate job interest: %s", resp.Status)
	}

	return nil
}
