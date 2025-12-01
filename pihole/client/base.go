package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type endpoint string

const (
	authHeader = "X-FTL-SID"
)

type Client struct {
	baseURL    string
	httpClient *http.Client

	password  string
	sessionID string
}

func NewClient(ctx context.Context, baseURL, password string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is empty")
	}
	// validate URL

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	c := &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
		password:   password,
	}
	sid, err := c.obtainSessionID(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not obtain session ID: %w", err)
	}
	c.sessionID = sid
	return c, nil
}

// get performs a GET request to the specified endpoint
func (c *Client) get(ctx context.Context, endpoint string) (*http.Response, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	if c.sessionID != "" {
		req.Header.Set(authHeader, c.sessionID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GET request: %w", err)
	}

	return resp, nil
}

// post performs a POST request to the specified endpoint with the given payload
func (c *Client) post(ctx context.Context, endpoint string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.sessionID != "" {
		req.Header.Set(authHeader, c.sessionID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute POST request: %w", err)
	}

	return resp, nil
}

// getJSON performs a GET request and decodes the JSON response into the target
func (c *Client) getJSON(ctx context.Context, endpoint string, target any) error {
	resp, err := c.get(ctx, endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return newHTTPError(resp.StatusCode, resp.Status, string(bodyBytes), http.MethodGet, c.baseURL+endpoint)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// postJSON performs a POST request and decodes the JSON response into the target
func (c *Client) postJSON(ctx context.Context, endpoint string, payload any, target any) error {
	resp, err := c.post(ctx, endpoint, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return newHTTPError(resp.StatusCode, resp.Status, string(bodyBytes), http.MethodPost, c.baseURL+endpoint)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
