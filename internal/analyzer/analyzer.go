package analyzer

import (
	"context"
	"time"

	"github.com/yannisobert/git-health-check/internal/github"
)

type Analyzer struct {
	client *github.Client
}

func New(client *github.Client) *Analyzer {
	return &Analyzer{client: client}
}

func (a *Analyzer) Analyze(ctx context.Context, owner, repo string) (*Report, error) {
	data, err := a.client.FetchAll(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	return buildReport(data, owner, repo), nil
}

func buildReport(data *github.RepoData, owner, repo string) *Report {
	categories := []Category{
		buildCategory("Maintenance", []CheckResult{
			checkREADME(data),
			checkLICENSE(data),
			checkCHANGELOG(data),
			checkCONTRIBUTING(data),
		}),
		buildCategory("Activity", []CheckResult{
			checkLastCommit(data),
			checkIssueRatio(data),
			checkPRMergeRatio(data),
		}),
		buildCategory("Conventions", []CheckResult{
			checkConventionalCommits(data),
			checkPRTemplate(data),
			checkSemverReleases(data),
			checkLinter(data),
		}),
		buildCategory("CI/CD", []CheckResult{
			checkGitHubActions(data),
			checkOtherCI(data),
			checkPreCommit(data),
			checkDependabot(data),
		}),
	}

	var score, max int
	for _, c := range categories {
		score += c.Score
		max += c.MaxScore
	}

	return &Report{
		Owner:      owner,
		Repo:       repo,
		FullName:   data.Repo.FullName,
		Score:      score,
		MaxScore:   max,
		Categories: categories,
		AnalyzedAt: time.Now(),
	}
}

func buildCategory(name string, checks []CheckResult) Category {
	var score, max int
	for _, c := range checks {
		score += c.Score
		max += c.MaxScore
	}
	return Category{Name: name, Score: score, MaxScore: max, Checks: checks}
}
