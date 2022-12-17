package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/db/client"
)

func TestNewHTTP(t *testing.T) {
	url := "http://localhost:54321"
	apikey := "apikey"
	q := client.NewHTTP(url, apikey)
	if q.URL != url {
		t.Errorf("want %v, got %v", url, q.URL)
	}
	if q.APIKey != apikey {
		t.Errorf("want %v, got %v", apikey, q.APIKey)
	}
}

// This will not compile if the interface is not implemented
func TestHTTPQueriesImplementsQuerier(t *testing.T) {
	var _ client.Querier = &client.HTTPQueries{}
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

	q := client.NewHTTP(ts.URL, apikey)

	ctx := context.Background()
	inputBytes, err := json.Marshal(input)
	if err != nil {
		t.Errorf("failed to marshal input: %v", err)
	}
	q.DoRequest(ctx, method, path, io.NopCloser(bytes.NewReader(
		inputBytes,
	)))
}
