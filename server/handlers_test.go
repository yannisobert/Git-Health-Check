package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	githubclient "github.com/yannisobert/git-health-check/internal/github"
)

// ── /api/health ──────────────────────────────────────────────────────────────

func TestHandleHealth_StatusAndBody(t *testing.T) {
	rr := httptest.NewRecorder()
	handleHealth(rr, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("status = %q, want ok", body["status"])
	}
	if body["app"] != "owlspector" {
		t.Errorf("app = %q, want owlspector", body["app"])
	}
}

// ── repoParam ─────────────────────────────────────────────────────────────────

func TestRepoParam_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?repo=owner/myrepo", nil)
	owner, repo, ok := repoParam(req, "repo")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if owner != "owner" || repo != "myrepo" {
		t.Errorf("got %q/%q, want owner/myrepo", owner, repo)
	}
}

func TestRepoParam_Invalid(t *testing.T) {
	cases := []string{"", "owneronly", "/repo", "owner/", "/"}
	for _, val := range cases {
		req := httptest.NewRequest(http.MethodGet, "/?repo="+val, nil)
		_, _, ok := repoParam(req, "repo")
		if ok {
			t.Errorf("expected ok=false for %q", val)
		}
	}
}

// ── writeAPIError ─────────────────────────────────────────────────────────────

func TestWriteAPIError_NotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	writeAPIError(rr, githubclient.ErrNotFound)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func TestWriteAPIError_RateLimit(t *testing.T) {
	rr := httptest.NewRecorder()
	writeAPIError(rr, githubclient.ErrRateLimit)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d, want 429", rr.Code)
	}
}

func TestWriteAPIError_Generic(t *testing.T) {
	rr := httptest.NewRecorder()
	writeAPIError(rr, errors.New("something broke"))
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
	var body map[string]string
	json.NewDecoder(rr.Body).Decode(&body)
	if body["error"] == "" {
		t.Error("expected non-empty error in body")
	}
}

// ── corsMiddleware ─────────────────────────────────────────────────────────────

func TestCORSMiddleware_SetsHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).
		ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected Access-Control-Allow-Origin: *")
	}
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("expected Access-Control-Allow-Methods to be set")
	}
}

func TestCORSMiddleware_Options(t *testing.T) {
	rr := httptest.NewRecorder()
	corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rr, httptest.NewRequest(http.MethodOptions, "/", nil))

	if rr.Code != http.StatusNoContent {
		t.Errorf("OPTIONS status = %d, want 204", rr.Code)
	}
}

// ── handlers: missing param → 400 ────────────────────────────────────────────

func TestAnalyzeHandler_MissingParam(t *testing.T) {
	rr := httptest.NewRecorder()
	(&handler{}).analyze(rr, httptest.NewRequest(http.MethodGet, "/api/analyze", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestHistoryHandler_MissingParam(t *testing.T) {
	rr := httptest.NewRecorder()
	(&handler{}).history(rr, httptest.NewRequest(http.MethodGet, "/api/history", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCompareHandler_MissingParam(t *testing.T) {
	rr := httptest.NewRecorder()
	(&handler{}).compare(rr, httptest.NewRequest(http.MethodGet, "/api/compare", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestRivalsHandler_MissingParam(t *testing.T) {
	rr := httptest.NewRecorder()
	(&handler{}).rivals(rr, httptest.NewRequest(http.MethodGet, "/api/rivals", nil))
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}
