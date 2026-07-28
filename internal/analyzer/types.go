package analyzer

import "time"

type Status string

const (
	StatusOK   Status = "OK"
	StatusWarn Status = "WARN"
	StatusFail Status = "FAIL"
)

type CheckResult struct {
	Name       string `json:"name"`
	Score      int    `json:"score"`
	MaxScore   int    `json:"maxScore"`
	Status     Status `json:"status"`
	Detail     string `json:"detail"`
	Suggestion string `json:"suggestion"`
}

type Category struct {
	Name     string        `json:"name"`
	Score    int           `json:"score"`
	MaxScore int           `json:"maxScore"`
	Checks   []CheckResult `json:"checks"`
}

type Report struct {
	Owner      string     `json:"owner"`
	Repo       string     `json:"repo"`
	FullName   string     `json:"fullName"`
	Score      int        `json:"score"`
	MaxScore   int        `json:"maxScore"`
	Categories []Category `json:"categories"`
	AnalyzedAt time.Time  `json:"analyzedAt"`
}

type HistoryPoint struct {
	Date  time.Time `json:"date"`
	Score int       `json:"score"`
	Event string    `json:"event,omitempty"`
}

// HistoryResult wraps the reconstructed score points with how much real
// history backs them, so callers can tell "this repo is just young" apart
// from "we stopped paging before reaching the requested span".
type HistoryResult struct {
	Points []HistoryPoint `json:"points"`
	// CoveredDays is the actual span, in days, between the oldest fetched
	// commit and now — it can be less than the period's target span either
	// because the repo isn't that old, or because Truncated is true.
	CoveredDays int `json:"coveredDays"`
	// Truncated is true when commit pagination hit its safety cap
	// (maxCommitPages) before reaching the requested span or the repo's
	// actual first commit — the true history may extend further back.
	Truncated bool `json:"truncated"`
}

type CompareResult struct {
	Repo1 Report               `json:"repo1"`
	Repo2 Report               `json:"repo2"`
	Diff  map[string]CheckDiff `json:"diff"`
}

type CheckDiff struct {
	Repo1Score int    `json:"repo1Score"`
	Repo2Score int    `json:"repo2Score"`
	Delta      int    `json:"delta"`
	Winner     string `json:"winner"`
}

type RivalSuggestion struct {
	FullName   string `json:"fullName"`
	Reason     string `json:"reason"`
	Similarity string `json:"similarity"`
}
