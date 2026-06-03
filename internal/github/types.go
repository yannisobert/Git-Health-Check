package github

import "time"

type Repo struct {
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Language        string    `json:"language"`
	DefaultBranch   string    `json:"default_branch"`
	PushedAt        time.Time `json:"pushed_at"`
	CreatedAt       time.Time `json:"created_at"`
	Topics          []string  `json:"topics"`
}

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type Issue struct {
	Number      int       `json:"number"`
	State       string    `json:"state"`
	PullRequest *struct{} `json:"pull_request,omitempty"`
}

func (i Issue) IsPR() bool { return i.PullRequest != nil }

type PullRequest struct {
	Number   int        `json:"number"`
	State    string     `json:"state"`
	MergedAt *time.Time `json:"merged_at"`
}

type Release struct {
	TagName     string    `json:"tag_name"`
	PublishedAt time.Time `json:"published_at"`
}

type ContentEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type Languages map[string]int

// RepoData bundles all data fetched in parallel for a repository.
type RepoData struct {
	Repo         Repo
	Commits      []Commit
	Issues       []Issue
	PullRequests []PullRequest
	Releases     []Release
	Languages    Languages
	RootContents []ContentEntry
	DotGithub    []ContentEntry
	Workflows    []ContentEntry
}
