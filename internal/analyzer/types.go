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
