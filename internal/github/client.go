package github

import (
	"context"
	"net/http"
	"os"
	"time"
)

const defaultBaseURL = "https://api.github.com"

// Client wraps GitHub REST API calls.
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// NewClient creates a GitHub API client using GITHUB_TOKEN from the environment.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      os.Getenv("GITHUB_TOKEN"),
		baseURL:    defaultBaseURL,
	}
}

// Ping verifies the client can reach the GitHub API.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/rate_limit", nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return ErrRateLimit
	}
	if resp.StatusCode != http.StatusOK {
		return ErrAPI{Status: resp.StatusCode}
	}
	return nil
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}
