package analyzer

import "github.com/yannisobert/git-health-check/internal/github"

// Analyzer orchestrates repository health checks.
type Analyzer struct {
	client *github.Client
}

func New(client *github.Client) *Analyzer {
	return &Analyzer{client: client}
}
