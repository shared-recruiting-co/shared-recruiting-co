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
	path := "/predict"
	input := "input"
	want := true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got %v, want %v", r.Method, "POST")
		}
		if r.URL.Path != path {
			t.Errorf("got %v, want %v", r.URL.Path, path)
		}
		if r.Header.Get("X-API-KEY") != apiKey {
			t.Errorf("got %v, want %v", r.Header.Get("X-API-KEY"), apiKey)
		}
		var body PredictRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if body.Input != input {
			t.Errorf("got %v, want %v", body.Input, input)
		}
		resp := PredictResponse{
			Result: want,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))

	client := NewClassifierClient(ctx, ClassifierClientArgs{
		BaseURL: ts.URL,
		ApiKey:  apiKey,
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
	path := "/predict/batch"
	inputs := map[string]string{
		"input1": "input1",
		"input2": "input2",
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
		if r.Header.Get("X-API-KEY") != apiKey {
			t.Errorf("got %v, want %v", r.Header.Get("X-API-KEY"), apiKey)
		}
		var body PredictBatchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		for k, v := range inputs {
			if body.Inputs[k] != v {
				t.Errorf("got %v, want %v", body.Inputs[k], v)
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
		BaseURL: ts.URL,
		ApiKey:  apiKey,
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
	input := "input"
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
