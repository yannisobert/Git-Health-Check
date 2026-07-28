package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) GetCommits(ctx context.Context, owner, repo string, limit int) ([]Commit, error) {
	var all []Commit
	page := 1
	for len(all) < limit {
		batch := 100
		if limit-len(all) < batch {
			batch = limit - len(all)
		}
		var commits []Commit
		url := fmt.Sprintf("/repos/%s/%s/commits?per_page=%d&page=%d", owner, repo, batch, page)
		if err := c.get(ctx, url, &commits); err != nil {
			return nil, err
		}
		all = append(all, commits...)
		if len(commits) < batch {
			break
		}
		page++
	}
	return all, nil
}

// GetCommitsPage returns one page of commits, newest first.
func (c *Client) GetCommitsPage(ctx context.Context, owner, repo string, perPage, page int) ([]Commit, error) {
	var commits []Commit
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/commits?per_page=%d&page=%d", owner, repo, perPage, page), &commits)
	return commits, err
}

// GetFirstCommitDate returns the date of the oldest commit touching path,
// within the most recent 100 commits affecting it. For files that are rarely
// modified (config files, templates) this is effectively the file's creation
// date; for heavily-edited files (e.g. README) it may understate how long the
// file has existed if it has been touched by more than 100 commits.
func (c *Client) GetFirstCommitDate(ctx context.Context, owner, repo, path string) (time.Time, error) {
	var commits []Commit
	err := c.get(ctx, fmt.Sprintf("/repos/%s/%s/commits?path=%s&per_page=100", owner, repo, url.QueryEscape(path)), &commits)
	if err != nil {
		return time.Time{}, err
	}
	if len(commits) == 0 {
		return time.Time{}, ErrNotFound
	}
	return commits[len(commits)-1].Commit.Author.Date, nil
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
