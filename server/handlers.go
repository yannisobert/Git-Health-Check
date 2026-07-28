package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/yannisobert/git-health-check/internal/analyzer"
	githubclient "github.com/yannisobert/git-health-check/internal/github"
)

type handler struct {
	analyzer *analyzer.Analyzer
}

func newHandler() *handler {
	return &handler{analyzer: analyzer.New(githubclient.NewClient())}
}

func (h *handler) analyze(w http.ResponseWriter, r *http.Request) {
	owner, repo, ok := repoParam(r, "repo")
	if !ok {
		writeError(w, http.StatusBadRequest, "missing or invalid ?repo=owner/repo")
		return
	}

	report, err := h.analyzer.Analyze(r.Context(), owner, repo)
	if err != nil {
		writeAPIError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, report)
}

func (h *handler) history(w http.ResponseWriter, r *http.Request) {
	owner, repo, ok := repoParam(r, "repo")
	if !ok {
		writeError(w, http.StatusBadRequest, "missing or invalid ?repo=owner/repo")
		return
	}

	period := analyzer.Period(r.URL.Query().Get("period"))
	if period != analyzer.PeriodMonthly {
		period = analyzer.PeriodWeekly
	}

	result, err := h.analyzer.History(r.Context(), owner, repo, period)
	if err != nil {
		writeAPIError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *handler) compare(w http.ResponseWriter, r *http.Request) {
	owner1, repo1, ok1 := repoParam(r, "repo1")
	owner2, repo2, ok2 := repoParam(r, "repo2")
	if !ok1 || !ok2 {
		writeError(w, http.StatusBadRequest, "missing or invalid ?repo1=owner/repo&repo2=owner/repo")
		return
	}

	result, err := h.analyzer.Compare(r.Context(), owner1, repo1, owner2, repo2)
	if err != nil {
		writeAPIError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *handler) rivals(w http.ResponseWriter, r *http.Request) {
	owner, repo, ok := repoParam(r, "repo")
	if !ok {
		writeError(w, http.StatusBadRequest, "missing or invalid ?repo=owner/repo")
		return
	}

	suggestions, err := h.analyzer.Rivals(r.Context(), owner, repo)
	if err != nil {
		writeAPIError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, suggestions)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "app": "owlspector"})
}

// helpers

func repoParam(r *http.Request, key string) (owner, repo string, ok bool) {
	val := r.URL.Query().Get(key)
	parts := strings.SplitN(val, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeAPIError(w http.ResponseWriter, err error) {
	if errors.Is(err, githubclient.ErrNotFound) {
		writeError(w, http.StatusNotFound, "repository not found")
		return
	}
	if errors.Is(err, githubclient.ErrRateLimit) {
		writeError(w, http.StatusTooManyRequests, "GitHub rate limit exceeded — add a GITHUB_TOKEN")
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}
