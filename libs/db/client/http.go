package client

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

// HTTPQueries is a client for the PostgREST API
type HTTPQueries struct {
	client *http.Client
	// URL is the Base URL of the PostgREST API
	URL string
	// APIKey is the API Key for the PostgREST API
	APIKey string
}

// NewHTTP creates a new HTTPQueries.
func NewHTTP(url, apiKey string) *HTTPQueries {
	return &HTTPQueries{
		client: &http.Client{},
		URL:    url,
		APIKey: apiKey,
	}
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

	req.Header.Set("apikey", q.APIKey)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", q.APIKey))
	req.Header.Set("Content-Type", "application/json")
	// always upsert on POST (aka insert)
	if method == http.MethodPost {
		req.Header.Set("Prefer", "resolution=merge-duplicates")
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
	basePath := "/user_profile"
	query := fmt.Sprintf("select=*&email=eq.%s", email)
	path := fmt.Sprintf("%s?%s", basePath, query)
	var result UserProfile

	resp, err := q.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("fetch user profile by email: %s", resp.Status)
	}

	var profile []UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return result, err
	}

	result, err = singleOrError(profile)
	return result, err
}

// GetUserEmailSyncHistory fetches a user's email sync history.
func (q *HTTPQueries) GetUserEmailSyncHistory(ctx context.Context, userID uuid.UUID) (UserEmailSyncHistory, error) {
	basePath := "/user_email_sync_history"
	query := fmt.Sprintf("select=*&user_id=eq.%s", userID)
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
	query := fmt.Sprintf("select=*&user_id=eq.%s&provider=eq.%s", arg.UserID, arg.Provider)
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
	path := fmt.Sprintf("%s", basePath)
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upsert user email sync history: %s", resp.Status)
	}

	return nil
}

// ListUserOAuthTokens lists a user oauth tokens.
// TODO: Support pagination
func (q *HTTPQueries) ListUserOAuthTokens(ctx context.Context, arg ListUserOAuthTokensParams) ([]UserOauthToken, error) {
	basePath := "/user_oauth_token"
	query := "select=*"
	if arg.Provider == "" {
		return nil, fmt.Errorf("provider is required")
	}
	query = fmt.Sprintf("%s&provider=eq.%s", query, arg.Provider)
	query = fmt.Sprintf("%s&is_valid=eq.%t", query, arg.IsValid)

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

	if result == nil || len(result) == 0 {
		// for now, return same error as database client
		return result, sql.ErrNoRows
	}

	return result, nil
}

// UpsertUserEmailSyncHistory upserts a user's email sync history.
func (q *HTTPQueries) UpsertUserEmailSyncHistory(ctx context.Context, arg UpsertUserEmailSyncHistoryParams) error {
	basePath := "/user_email_sync_history"
	path := fmt.Sprintf("%s", basePath)
	body, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	resp, err := q.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upsert user email sync history: %s", resp.Status)
	}

	return nil
}
