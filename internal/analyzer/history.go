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

// History returns the score evolution over time for a repository.
func (a *Analyzer) History(ctx context.Context, owner, repo string, period Period) ([]HistoryPoint, error) {
	data, err := a.client.FetchAll(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	// Fetch more commits for richer history
	commits, err := a.client.GetCommits(ctx, owner, repo, 100)
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return nil, nil
	}
	data.Commits = commits

	report := buildReport(data, owner, repo)

	// Split score into fixed (non-activity) and activity
	var fixedScore, activityMax int
	for _, cat := range report.Categories {
		if cat.Name == "Activity" {
			activityMax = cat.MaxScore
		} else {
			fixedScore += cat.Score
		}
	}

	return buildHistory(commits, period, fixedScore, activityMax), nil
}

func buildHistory(commits []github.Commit, period Period, fixedScore, activityMax int) []HistoryPoint {
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
			Score: fixedScore + actScore,
			Event: event,
		})

		start = end
	}

	if len(points) > maxPoints {
		points = points[len(points)-maxPoints:]
	}
	return points
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
