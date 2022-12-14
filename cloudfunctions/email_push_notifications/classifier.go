package cloudfunctions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Consider moving to share library in future

type PredictRequest struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type PredictResponse struct {
	Result bool `json:"result"`
}

type PredictBatchRequest struct {
	Inputs map[string]*PredictRequest `json:"inputs"`
}

type PredictBatchResponse struct {
	Results map[string]bool `json:"results"`
}

type Classifier interface {
	// Predict if the input meets the classification
	Predict(input *PredictRequest) (bool, error)
	// Predict if the inputs are spam or not
	PredictBatch(inputs map[string]*PredictRequest) (map[string]bool, error)
}

type ClassifierClient struct {
	ctx       context.Context
	baseURL   string
	authToken string
}

type ClassifierClientArgs struct {
	BaseURL   string
	ApiKey    string
	AuthToken string
}

func NewClassifierClient(ctx context.Context, args ClassifierClientArgs) *ClassifierClient {
	return &ClassifierClient{
		ctx:       ctx,
		baseURL:   args.BaseURL,
		authToken: args.AuthToken,
	}
}

func (c *ClassifierClient) Predict(input *PredictRequest) (bool, error) {
	resp := &PredictResponse{}
	err := c.doRequest("POST", "/v1/predict", input, resp)
	if err != nil {
		return false, err
	}
	return resp.Result, nil
}

func (c *ClassifierClient) PredictBatch(inputs map[string]*PredictRequest) (map[string]bool, error) {
	req := &PredictBatchRequest{Inputs: inputs}
	resp := &PredictBatchResponse{}
	err := c.doRequest("POST", "/v1/predict/batch", req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (c *ClassifierClient) doRequest(method string, path string, req interface{}, resp interface{}) error {
	url := c.baseURL + path
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(body)
	httpReq, err := http.NewRequestWithContext(c.ctx, method, url, reqBody)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.authToken != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	}
	httpResp, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return err
	}

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return err
}
