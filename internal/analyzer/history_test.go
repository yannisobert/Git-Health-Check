package analyzer

import (
	"testing"
	"time"

	"github.com/yannisobert/git-health-check/internal/github"
)

func TestTruncateToPeriod_Monthly(t *testing.T) {
	d := time.Date(2024, 3, 15, 12, 30, 0, 0, time.UTC)
	got := truncateToPeriod(d, PeriodMonthly)
	want := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTruncateToPeriod_Weekly_Wednesday(t *testing.T) {
	// Wednesday 2024-03-13 → Monday 2024-03-11
	d := time.Date(2024, 3, 13, 10, 0, 0, 0, time.UTC)
	got := truncateToPeriod(d, PeriodWeekly)
	want := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTruncateToPeriod_Weekly_Sunday(t *testing.T) {
	// Sunday 2024-03-17 → Monday 2024-03-11
	d := time.Date(2024, 3, 17, 10, 0, 0, 0, time.UTC)
	got := truncateToPeriod(d, PeriodWeekly)
	want := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestTruncateToPeriod_Weekly_Monday(t *testing.T) {
	// Monday stays unchanged
	d := time.Date(2024, 3, 11, 8, 0, 0, 0, time.UTC)
	got := truncateToPeriod(d, PeriodWeekly)
	want := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAddPeriod_Weekly(t *testing.T) {
	d := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC)
	got := addPeriod(d, PeriodWeekly)
	want := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAddPeriod_Monthly(t *testing.T) {
	d := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	got := addPeriod(d, PeriodMonthly)
	want := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestActivityScoreForDays(t *testing.T) {
	tests := []struct {
		days int
		max  int
		want int
	}{
		{0, 10, 10},
		{29, 10, 10},
		{30, 10, 7},
		{89, 10, 7},
		{90, 10, 4},
		{179, 10, 4},
		{180, 10, 0},
		{9999, 10, 0},
	}
	for _, tt := range tests {
		got := activityScoreForDays(tt.days, tt.max)
		if got != tt.want {
			t.Errorf("activityScoreForDays(%d, %d) = %d, want %d", tt.days, tt.max, got, tt.want)
		}
	}
}

func TestDetectEvent_InRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	commits := []github.Commit{commit("feat: add ci workflow", 0)}
	commits[0].Commit.Author.Date = time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)

	got := detectEvent(commits, start, end)
	if got != "CI added" {
		t.Errorf("got %q, want CI added", got)
	}
}

func TestDetectEvent_OutOfRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	commits := []github.Commit{commit("feat: add ci workflow", 0)}
	commits[0].Commit.Author.Date = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	got := detectEvent(commits, start, end)
	if got != "" {
		t.Errorf("expected no event, got %q", got)
	}
}

func TestDetectEvent_NoKeyword(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	commits := []github.Commit{commit("feat: add button component", 0)}
	commits[0].Commit.Author.Date = time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)

	got := detectEvent(commits, start, end)
	if got != "" {
		t.Errorf("expected no event, got %q", got)
	}
}

func TestBuildHistory_Empty(t *testing.T) {
	if buildHistory(nil, PeriodWeekly, 50, 25) != nil {
		t.Error("expected nil for empty commits")
	}
}

func TestBuildHistory_ScoresInRange(t *testing.T) {
	commits := []github.Commit{
		commit("feat: a", 1),
		commit("fix: b", 8),
		commit("chore: c", 15),
		commit("docs: d", 22),
	}
	maxTotal := 75
	points := buildHistory(commits, PeriodWeekly, 50, 25)
	if len(points) == 0 {
		t.Fatal("expected at least one history point")
	}
	for _, p := range points {
		if p.Score < 0 || p.Score > maxTotal {
			t.Errorf("score %d out of range [0, %d]", p.Score, maxTotal)
		}
	}
}

func TestBuildHistory_Monthly_CappedAt12(t *testing.T) {
	// 14 months of commits (newest first)
	var commits []github.Commit
	for i := 0; i < 14; i++ {
		commits = append(commits, commit("feat: x", i*31))
	}
	points := buildHistory(commits, PeriodMonthly, 50, 25)
	if len(points) > 12 {
		t.Errorf("got %d points, want at most 12", len(points))
	}
}

func TestBuildHistory_Weekly_CappedAt26(t *testing.T) {
	// 30 weeks of commits
	var commits []github.Commit
	for i := 0; i < 30; i++ {
		commits = append(commits, commit("feat: x", i*7))
	}
	points := buildHistory(commits, PeriodWeekly, 50, 25)
	if len(points) > 26 {
		t.Errorf("got %d points, want at most 26", len(points))
	}
}
