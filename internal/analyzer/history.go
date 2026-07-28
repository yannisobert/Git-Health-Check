package analyzer

import (
	"context"
	"strings"
	"time"

	"github.com/yannisobert/git-health-check/internal/github"
)

type Period string

const (
	PeriodWeekly  Period = "weekly"
	PeriodMonthly Period = "monthly"
)

// Target lookback per period so "Weekly" and "1Y" reflect a real span of time
// rather than however far back a fixed commit count happens to reach.
const (
	weeklySpan  = 26 * 7 * 24 * time.Hour // ~6 months, matches the 26-point weekly cap
	monthlySpan = 365 * 24 * time.Hour    // 1 year, matches the 12-point monthly cap

	// maxCommitPages caps pagination at 1000 commits so a single history
	// request can't run away against repos with enormous commit counts.
	maxCommitPages = 10
	commitsPerPage = 100
)

func targetSpan(period Period) time.Duration {
	if period == PeriodMonthly {
		return monthlySpan
	}
	return weeklySpan
}

// fileTimelines holds the first-appearance date of each file or directory
// that backs a presence check, resolved once so buildHistory can cheaply
// compare against it at every period point instead of re-querying GitHub per
// point. A zero value means the file was never found.
type fileTimelines struct {
	readme, license, changelog, contributing  time.Time
	prTemplate, linter                        time.Time
	ghActions, otherCI, preCommit, dependabot time.Time
}

// History returns the score evolution over time for a repository, along
// with how much real history actually backs it (see HistoryResult).
func (a *Analyzer) History(ctx context.Context, owner, repo string, period Period) (*HistoryResult, error) {
	data, err := a.client.FetchAll(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	// Paginate back far enough to give this period a real span of history,
	// instead of stopping at a fixed commit count regardless of repo velocity.
	commits, truncated, err := a.fetchCommitsSince(ctx, owner, repo, time.Now().Add(-targetSpan(period)))
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return &HistoryResult{Points: []HistoryPoint{}}, nil
	}
	data.Commits = commits

	oldest := commits[len(commits)-1].Commit.Author.Date
	timelines := a.resolveFileTimelines(ctx, owner, repo, data, oldest)

	report := buildReport(data, owner, repo)
	var activityMax int
	for _, cat := range report.Categories {
		if cat.Name == "Activity" {
			activityMax = cat.MaxScore
		}
	}

	return &HistoryResult{
		Points:      buildHistory(data, period, timelines, activityMax),
		CoveredDays: int(time.Since(oldest).Hours() / 24),
		Truncated:   truncated,
	}, nil
}

// fetchCommitsSince pages through commits (newest first) until the oldest
// commit on a page predates cutoff, the repo runs out of commits, or
// maxCommitPages is reached — whichever comes first. The returned bool is
// true only in the last case: pagination was cut off by the safety cap
// before reaching cutoff or the repo's actual first commit.
func (a *Analyzer) fetchCommitsSince(ctx context.Context, owner, repo string, cutoff time.Time) ([]github.Commit, bool, error) {
	var all []github.Commit
	for page := 1; page <= maxCommitPages; page++ {
		batch, err := a.client.GetCommitsPage(ctx, owner, repo, commitsPerPage, page)
		if err != nil {
			return nil, false, err
		}
		if len(batch) == 0 {
			return all, false, nil
		}
		all = append(all, batch...)
		if batch[len(batch)-1].Commit.Author.Date.Before(cutoff) {
			return all, false, nil
		}
		if len(batch) < commitsPerPage {
			return all, false, nil
		}
	}
	return all, true, nil
}

// resolveFileTimelines finds, once per request, the date each file or
// directory backing a presence check first appeared. Only files GitHub
// currently reports as present are queried. If a lookup fails (rate limit,
// transient error), the file is assumed to have existed since the oldest
// known commit, so that single check degrades to today's frozen behavior
// instead of failing the whole history request.
func (a *Analyzer) resolveFileTimelines(ctx context.Context, owner, repo string, data *github.RepoData, fallback time.Time) fileTimelines {
	var t fileTimelines

	resolve := func(path string) time.Time {
		date, err := a.client.GetFirstCommitDate(ctx, owner, repo, path)
		if err != nil {
			return fallback
		}
		return date
	}

	if e, ok := findFile(data.RootContents, "readme.md", "readme", "readme.txt", "readme.rst"); ok {
		t.readme = resolve(e.Path)
	}
	if e, ok := findFile(data.RootContents, "license", "license.md", "license.txt", "licence", "copying"); ok {
		t.license = resolve(e.Path)
	}
	if e, ok := findFile(data.RootContents, "changelog.md", "changelog", "changelog.txt", "history.md"); ok {
		t.changelog = resolve(e.Path)
	}
	if e, ok := findFile(data.RootContents, "contributing.md", "contributing", "contributing.txt"); ok {
		t.contributing = resolve(e.Path)
	}
	if e, ok := findFile(data.DotGithub, "pull_request_template.md"); ok {
		t.prTemplate = resolve(e.Path)
	}
	if e, ok := findFile(data.RootContents, linterFiles...); ok {
		t.linter = resolve(e.Path)
	}
	if len(data.Workflows) > 0 {
		t.ghActions = resolve(".github/workflows")
	}
	if e, ok := findFile(data.RootContents, otherCIFiles...); ok {
		t.otherCI = resolve(e.Path)
	}
	if e, ok := findFile(data.RootContents, preCommitFiles...); ok {
		t.preCommit = resolve(e.Path)
	}
	if e, ok := findFile(data.DotGithub, "dependabot.yml", "dependabot.yaml"); ok {
		t.dependabot = resolve(e.Path)
	}

	return t
}

func buildHistory(data *github.RepoData, period Period, t fileTimelines, activityMax int) []HistoryPoint {
	commits := data.Commits
	if len(commits) == 0 {
		return nil
	}

	maxPoints := 26
	if period == PeriodMonthly {
		maxPoints = 12
	}

	now := time.Now()
	oldest := commits[len(commits)-1].Commit.Author.Date
	start := truncateToPeriod(oldest, period)

	var points []HistoryPoint
	for !start.After(now) {
		end := addPeriod(start, period)

		// Most recent commit before end of this period (commits are newest-first)
		lastCommitDays := 9999
		for _, c := range commits {
			if c.Commit.Author.Date.Before(end) {
				lastCommitDays = int(end.Sub(c.Commit.Author.Date).Hours() / 24)
				break
			}
		}

		actScore := activityScoreForDays(lastCommitDays, activityMax)
		event := detectEvent(commits, start, end)

		points = append(points, HistoryPoint{
			Date:  start,
			Score: scoreAt(data, t, end) + actScore,
			Event: event,
		})

		start = end
	}

	if len(points) > maxPoints {
		points = points[len(points)-maxPoints:]
	}
	return points
}

// scoreAt reconstructs the Maintenance + Conventions + CI/CD score as of end,
// using each file's real first-appearance date and by filtering timestamped
// data (commits, releases) instead of freezing every check at its current
// value. Size/quality nuances (e.g. a short README) are not reconstructable
// without fetching historical file contents, so a present file counts for
// full points here even if today's live check only awards partial credit.
func scoreAt(data *github.RepoData, t fileTimelines, end time.Time) int {
	present := func(at time.Time) int {
		if at.IsZero() || at.After(end) {
			return 0
		}
		return 1
	}

	score := 0

	// Maintenance (30)
	score += present(t.readme) * 10
	score += present(t.license) * 10
	score += present(t.changelog) * 5
	score += present(t.contributing) * 5

	// Conventions (25)
	score += conventionalCommitsScoreAt(data.Commits, end)
	score += present(t.prTemplate) * 5
	score += semverReleasesScoreAt(data.Releases, end)
	score += present(t.linter) * 5

	// CI/CD (20)
	score += present(t.ghActions) * 8
	score += present(t.otherCI) * 4
	score += present(t.preCommit) * 4
	score += present(t.dependabot) * 4

	return score
}

func conventionalCommitsScoreAt(commits []github.Commit, end time.Time) int {
	var total, matching int
	for _, c := range commits {
		if !c.Commit.Author.Date.Before(end) {
			continue
		}
		total++
		msg := strings.ToLower(c.Commit.Message)
		for _, prefix := range conventionalPrefixes {
			if strings.HasPrefix(msg, prefix) {
				matching++
				break
			}
		}
	}
	if total == 0 {
		return 0
	}
	ratio := float64(matching) / float64(total)
	switch {
	case ratio >= 0.8:
		return 10
	case ratio >= 0.5:
		return 6
	case ratio >= 0.2:
		return 3
	default:
		return 0
	}
}

func semverReleasesScoreAt(releases []github.Release, end time.Time) int {
	for _, r := range releases {
		if r.PublishedAt.Before(end) && semverRe.MatchString(r.TagName) {
			return 5
		}
	}
	return 0
}

func activityScoreForDays(days, max int) int {
	switch {
	case days < 30:
		return max
	case days < 90:
		return max * 7 / 10
	case days < 180:
		return max * 4 / 10
	default:
		return 0
	}
}

func detectEvent(commits []github.Commit, start, end time.Time) string {
	keywords := map[string]string{
		"license":      "LICENSE added",
		"ci":           "CI added",
		"workflow":     "CI added",
		"dependabot":   "Dependabot added",
		"contributing": "CONTRIBUTING added",
		"changelog":    "CHANGELOG added",
		"release":      "Release",
	}
	for _, c := range commits {
		d := c.Commit.Author.Date
		if d.Before(start) || !d.Before(end) {
			continue
		}
		msg := strings.ToLower(c.Commit.Message)
		for kw, label := range keywords {
			if strings.Contains(msg, kw) {
				return label
			}
		}
	}
	return ""
}

func truncateToPeriod(t time.Time, period Period) time.Time {
	if period == PeriodMonthly {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	}
	// Weekly: truncate to Monday
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := t.AddDate(0, 0, -(weekday - 1))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
}

func addPeriod(t time.Time, period Period) time.Time {
	if period == PeriodMonthly {
		return t.AddDate(0, 1, 0)
	}
	return t.AddDate(0, 0, 7)
}
