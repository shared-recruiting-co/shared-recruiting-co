package cloudfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func TestClassifierClientPredict(t *testing.T) {
	ctx := context.Background()
	apiKey := "test"
	authToken := "xxx.yyy.zzz"
	path := "/v1/predict"
	input := &PredictRequest{
		From:    "from",
		Subject: "subject",
		Body:    "body",
	}
	want := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got %v, want %v", r.Method, "POST")
		}
		if r.URL.Path != path {
			t.Errorf("got %v, want %v", r.URL.Path, path)
		}
		wantAuth := fmt.Sprintf("Bearer %s", authToken)
		if r.Header.Get("Authorization") != wantAuth {
			t.Errorf("got %v, want %v", r.Header.Get("Authorization"), wantAuth)
		}
		var body PredictRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if body.From != input.From {
			t.Errorf("got %v, want %v", body.From, input.From)
		}
		if body.Subject != input.Subject {
			t.Errorf("got %v, want %v", body.Subject, input.Subject)
		}
		if body.Body != input.Body {
			t.Errorf("got %v, want %v", body.Body, input.Body)
		}
		resp := PredictResponse{
			Result: want,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))

	client := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL:   ts.URL,
		ApiKey:    apiKey,
		AuthToken: authToken,
	})

	got, err := client.Predict(input)

	if err != nil {
		t.Errorf("failed to predict: %v", err)
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestClassifierClientPredictBatch(t *testing.T) {
	ctx := context.Background()
	apiKey := "test"
	authToken := "xxx.yyy.zzz"
	path := "/v1/predict/batch"
	inputs := map[string]*PredictRequest{
		"input1": &PredictRequest{
			From:    "1",
			Subject: "1",
			Body:    "1",
		},
		"input2": &PredictRequest{
			From:    "2",
			Subject: "2",
			Body:    "2",
		},
	}
	want := map[string]bool{
		"input1": true,
		"input2": false,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got %v, want %v", r.Method, "POST")
		}
		if r.URL.Path != path {
			t.Errorf("got %v, want %v", r.URL.Path, path)
		}
		wantAuth := fmt.Sprintf("Bearer %s", authToken)
		if r.Header.Get("Authorization") != wantAuth {
			t.Errorf("got %v, want %v", r.Header.Get("Authorization"), wantAuth)
		}
		var body PredictBatchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		for k, v := range inputs {
			if body.Inputs[k].From != v.From {
				t.Errorf("got %v, want %v", body.Inputs[k].From, v.From)
			}
			if body.Inputs[k].Subject != v.Subject {
				t.Errorf("got %v, want %v", body.Inputs[k].Subject, v.Subject)
			}
			if body.Inputs[k].Body != v.Body {
				t.Errorf("got %v, want %v", body.Inputs[k].Body, v.Body)
			}
		}
		resp := PredictBatchResponse{
			Results: want,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))

	client := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL:   ts.URL,
		ApiKey:    apiKey,
		AuthToken: authToken,
	})

	got, err := client.PredictBatch(inputs)

	if err != nil {
		t.Errorf("failed to predict: %v", err)
	}

	for k, v := range got {
		if v != want[k] {
			t.Errorf("got %v, want %v", v, want[k])
		}
	}

}

func TestClassifierClientNon200(t *testing.T) {
	ctx := context.Background()
	apiKey := "test"
	input := &PredictRequest{
		From:    "from",
		Subject: "subject",
		Body:    "body",
	}
	want := false
	status := http.StatusBadRequest

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))

	client := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL: ts.URL,
		ApiKey:  apiKey,
	})

	got, err := client.Predict(input)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("%d", status)) {
		t.Errorf("expected %d to be in error message: %v", status, err)
	}
}
