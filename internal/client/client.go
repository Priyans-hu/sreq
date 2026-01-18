package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Priyans-hu/sreq/pkg/types"
)

// Client is the HTTP client for making requests
type Client struct {
	httpClient *http.Client
	verbose    bool
}

// New creates a new HTTP client
func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Option is a function that configures the client
type Option func(*Client)

// WithTimeout sets the HTTP client timeout
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithVerbose enables verbose output
func WithVerbose(v bool) Option {
	return func(c *Client) {
		c.verbose = v
	}
}

// Do executes an HTTP request
func (c *Client) Do(ctx context.Context, req *types.Request, creds *types.ResolvedCredentials) (*types.Response, error) {
	// Build the full URL
	url := fmt.Sprintf("%s%s", creds.BaseURL, req.Path)

	// Create the HTTP request
	var body io.Reader
	if req.Body != "" {
		body = bytes.NewBufferString(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set auth if available
	if creds.Username != "" && creds.Password != "" {
		httpReq.SetBasicAuth(creds.Username, creds.Password)
	}

	// Set additional credential headers
	for key, value := range creds.Headers {
		httpReq.Header.Set(key, value)
	}

	// Set default content type for requests with body
	if body != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	if c.verbose {
		fmt.Printf("> %s %s\n", req.Method, url)
		for key, values := range httpReq.Header {
			for _, value := range values {
				fmt.Printf("> %s: %s\n", key, value)
			}
		}
		fmt.Println(">")
	}

	// Execute the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &types.Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       respBody,
	}, nil
}
