package ml

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type EmailInput struct {
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type ClassifyRequest = EmailInput

type ClassifyResponse struct {
	Result bool `json:"result"`
}

type BatchClassifyRequest struct {
	Inputs map[string]*ClassifyRequest `json:"inputs"`
}

type BatchClassifyResponse struct {
	Results map[string]bool `json:"results"`
}

type ParseJobRequest = EmailInput

type ParseJobResponse struct {
	Company   string `json:"company_name"`
	Title     string `json:"job_title"`
	Recruiter string `json:"recruiter_name"`
}

type BatchParseJobRequest struct {
	Inputs map[string]*BatchParseJobRequest `json:"inputs"`
}

type BatchParseJobResponse struct {
	Results map[string]*ParseJobResponse `json:"results"`
}

type Service interface {
	// Classify if the input is a recruiting email or not
	Classify(input *ClassifyRequest) (*ClassifyResponse, error)
	// Classify if the inputs recruiting emails or not
	BatchClassify(inputs *BatchClassifyRequest) (*BatchClassifyResponse, error)
	// Parse if the input meets the classification
	ParseJob(input *ParseJobRequest) (*ParseJobResponse, error)
	// Parse if the inputs are spam or not
	BatchParseJob(req *BatchParseJobRequest) (*BatchParseJobResponse, error)
}

type client struct {
	ctx       context.Context
	baseURL   string
	authToken string
}

type NewServiceArg struct {
	BaseURL   string
	AuthToken string
}

func NewService(ctx context.Context, arg NewServiceArg) *client {
	return &client{
		ctx:       ctx,
		baseURL:   arg.BaseURL,
		authToken: arg.AuthToken,
	}
}

func (c *client) Classify(input *ClassifyRequest) (*ClassifyResponse, error) {
	resp := &ClassifyResponse{}
	err := c.doRequest("POST", "/v1/classify", input, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) BatchClassify(inputs *BatchClassifyRequest) (*BatchClassifyResponse, error) {
	resp := &BatchClassifyResponse{}
	err := c.doRequest("POST", "/v1/classify/batch", inputs, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) ParseJob(input *ParseJobRequest) (*ParseJobResponse, error) {
	resp := &ParseJobResponse{}
	err := c.doRequest("POST", "/v1/parse", input, resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (c *client) BatchParseJob(inputs *BatchParseJobRequest) (*BatchParseJobResponse, error) {
	resp := &BatchParseJobResponse{}
	err := c.doRequest("POST", "/v1/parse/batch", inputs, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) doRequest(method string, endpoint string, req interface{}, resp interface{}) error {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	reqURL := baseURL.ResolveReference(&url.URL{Path: path.Join(baseURL.Path, endpoint)})
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(body)
	httpReq, err := http.NewRequestWithContext(c.ctx, method, reqURL.String(), reqBody)
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
