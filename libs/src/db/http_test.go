package db_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/db"
)

func TestNewHTTP(t *testing.T) {
	url := "http://localhost:54321"
	apikey := "apikey"
	q := db.NewHTTP(url, apikey)
	if q.URL != url {
		t.Errorf("want %v, got %v", url, q.URL)
	}
	if q.APIKey != apikey {
		t.Errorf("want %v, got %v", apikey, q.APIKey)
	}
}

// This will not compile if the interface is not implemented
func TestHTTPQueriesImplementsQuerier(t *testing.T) {
	var _ db.Querier = &db.HTTPQueries{}
}

// Test DoRequest
func TestHTTPQueriesDoRequest(t *testing.T) {
	apikey := "apikey"
	path := "/test"
	method := http.MethodPost

	input := map[string]string{
		"test": "test",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Errorf("got %v, want %v", r.Method, method)
		}
		if r.URL.Path != path {
			t.Errorf("got %v, want %v", r.URL.Path, path)
		}
		// check auth headers
		wantAuth := fmt.Sprintf("Bearer %s", apikey)
		if r.Header.Get("Authorization") != wantAuth {
			t.Errorf("got %v, want %v", r.Header.Get("Authorization"), wantAuth)
		}
		if r.Header.Get("apikey") != apikey {
			t.Errorf("got %v, want %v", r.Header.Get("apikey"), apikey)
		}
		// check for upsert header
		if r.Method == http.MethodPost && r.Header.Get("Prefer") != "resolution=merge-duplicates" {
			t.Errorf("got %v, want %v", r.Header.Get("Prefer"), "resolution=merge-duplicates")
		}
		// check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("got %v, want %v", r.Header.Get("Content-Type"), "application/json")
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if body["test"] != input["test"] {
			t.Errorf("got %v, want %v", body["test"], input["test"])
		}
	}))
	defer ts.Close()

	q := db.NewHTTP(ts.URL, apikey)

	ctx := context.Background()
	inputBytes, err := json.Marshal(input)
	if err != nil {
		t.Errorf("failed to marshal input: %v", err)
	}
	_, _ = q.DoRequest(ctx, method, path, io.NopCloser(bytes.NewReader(
		inputBytes,
	)))
}

func TestHTTPGetUserProfileByEmail(t *testing.T) {
	apikey := "apikey"
	email := "example@test.com"
	wantPath := fmt.Sprintf("/user_profile?select=*&email=eq.%s", email)
	want := db.UserProfile{
		UserID:    uuid.New(),
		Email:     email,
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("got %v, want %v", r.Method, http.MethodGet)
		}
		if r.URL.String() != wantPath {
			t.Errorf("got %v, want %v", r.URL.String(), wantPath)
		}
		// check auth headers
		wantAuth := fmt.Sprintf("Bearer %s", apikey)
		if r.Header.Get("Authorization") != wantAuth {
			t.Errorf("got %v, want %v", r.Header.Get("Authorization"), wantAuth)
		}
		if r.Header.Get("apikey") != apikey {
			t.Errorf("got %v, want %v", r.Header.Get("apikey"), apikey)
		}
		// check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("got %v, want %v", r.Header.Get("Content-Type"), "application/json")
		}

		// return a dummy list of user profiles
		resp := []db.UserProfile{want}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	q := db.NewHTTP(ts.URL, apikey)

	ctx := context.Background()
	resp, err := q.GetUserProfileByEmail(ctx, email)
	if err != nil {
		t.Errorf("failed to get user profile: %v", err)
	}

	if resp.UserID.String() != want.UserID.String() {
		t.Errorf("got %v, want %v", resp.UserID, want.UserID)
	}
	if resp.Email != want.Email {
		t.Errorf("got %v, want %v", resp.Email, want.Email)
	}
	if resp.FirstName != want.FirstName {
		t.Errorf("got %v, want %v", resp.FirstName, want.FirstName)
	}
	if resp.LastName != want.LastName {
		t.Errorf("got %v, want %v", resp.LastName, want.LastName)
	}
}
