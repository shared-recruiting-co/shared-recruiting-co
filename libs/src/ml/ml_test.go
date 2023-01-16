package ml_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	ml "github.com/shared-recruiting-co/shared-recruiting-co/libs/src/ml"
)

func TestMLServiceClassify(t *testing.T) {
	ctx := context.Background()
	authToken := "xxx.yyy.zzz"
	path := "/v1/classify"
	input := &ml.ClassifyRequest{
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
		var body ml.ClassifyRequest
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
		resp := ml.ClassifyResponse{
			Result: want,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))

	srv := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   ts.URL,
		AuthToken: authToken,
	})

	got, err := srv.Classify(input)

	if err != nil {
		t.Errorf("failed to predict: %v", err)
	}
	if got.Result != want {
		t.Errorf("got %v, want %v", got.Result, want)
	}
}

func TestClassifierClientPredictBatch(t *testing.T) {
	ctx := context.Background()
	authToken := "xxx.yyy.zzz"
	path := "/v1/classify/batch"
	inputs := ml.BatchClassifyRequest{
		Inputs: map[string]*ml.ClassifyRequest{
			"input1": &ml.ClassifyRequest{
				From:    "1",
				Subject: "1",
				Body:    "1",
			},
			"input2": &ml.ClassifyRequest{
				From:    "2",
				Subject: "2",
				Body:    "2",
			},
		},
	}
	want := ml.BatchClassifyResponse{
		Results: map[string]bool{
			"input1": true,
			"input2": false,
		},
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
		var body ml.BatchClassifyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		for k, v := range inputs.Inputs {
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
		if err := json.NewEncoder(w).Encode(want); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))

	srv := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL:   ts.URL,
		AuthToken: authToken,
	})

	got, err := srv.BatchClassify(&inputs)

	if err != nil {
		t.Errorf("failed to predict: %v", err)
	}

	for k, v := range got.Results {
		if v != want.Results[k] {
			t.Errorf("got %v, want %v", v, want.Results[k])
		}
	}

}

func TestClassifierClientNon200(t *testing.T) {
	ctx := context.Background()
	input := &ml.ClassifyRequest{
		From:    "from",
		Subject: "subject",
		Body:    "body",
	}
	status := http.StatusBadRequest

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))

	client := ml.NewService(ctx, ml.NewServiceArg{
		BaseURL: ts.URL,
	})

	got, err := client.Classify(input)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if got != nil {
		t.Errorf("expected nil response, got %v", got)
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("%d", status)) {
		t.Errorf("expected %d to be in error message: %v", status, err)
	}
}
