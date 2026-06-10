package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const defaultBaseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      os.Getenv("GITHUB_TOKEN"),
		baseURL:    defaultBaseURL,
	}
}

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

func (c *Client) GetRepo(ctx context.Context, owner, repo string) (Repo, error) {
	var r Repo
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s", owner, repo), &r)
	return r, err
}

func (c *Client) GetCommits(ctx context.Context, owner, repo string, perPage int) ([]Commit, error) {
	var commits []Commit
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/commits?per_page=%d", owner, repo, perPage), &commits)
	return commits, err
}

func (c *Client) GetIssues(ctx context.Context, owner, repo string) ([]Issue, error) {
	var issues []Issue
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/issues?state=all&per_page=100", owner, repo), &issues)
	return issues, err
}

func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) ([]PullRequest, error) {
	var prs []PullRequest
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/pulls?state=all&per_page=100", owner, repo), &prs)
	return prs, err
}

func (c *Client) GetReleases(ctx context.Context, owner, repo string) ([]Release, error) {
	var releases []Release
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/releases?per_page=30", owner, repo), &releases)
	return releases, err
}

func (c *Client) GetLanguages(ctx context.Context, owner, repo string) (Languages, error) {
	langs := make(Languages)
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/languages", owner, repo), &langs)
	return langs, err
}

// GetDirContents returns the entries of a directory path. Returns nil if path does not exist.
func (c *Client) GetDirContents(ctx context.Context, owner, repo, path string) ([]ContentEntry, error) {
	var entries []ContentEntry
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path), &entries)
	if errors.Is(err, ErrNotFound) {
		return nil, nil
	}
	return entries, err
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return json.NewDecoder(resp.Body).Decode(out)
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusForbidden, http.StatusTooManyRequests:
		return ErrRateLimit
	default:
		return ErrAPI{Status: resp.StatusCode}
	}
}
