package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(srv *httptest.Server) *Client {
	return &Client{
		httpClient: srv.Client(),
		baseURL:    srv.URL,
	}
}

func TestGet_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	var out Repo
	err := newTestClient(srv).get(context.Background(), "/repos/foo/bar", &out)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestGet_429(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	var out Repo
	err := newTestClient(srv).get(context.Background(), "/repos/foo/bar", &out)
	if !errors.Is(err, ErrRateLimit) {
		t.Errorf("want ErrRateLimit, got %v", err)
	}
}

func TestGet_403(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()
	var out Repo
	err := newTestClient(srv).get(context.Background(), "/repos/foo/bar", &out)
	if !errors.Is(err, ErrRateLimit) {
		t.Errorf("want ErrRateLimit, got %v", err)
	}
}

func TestGet_5xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	var out Repo
	err := newTestClient(srv).get(context.Background(), "/repos/foo/bar", &out)
	var apiErr ErrAPI
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusInternalServerError {
		t.Errorf("want ErrAPI{500}, got %v", err)
	}
}

func TestGet_200_Decodes(t *testing.T) {
	want := Repo{FullName: "foo/bar", Language: "Go"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer srv.Close()
	var got Repo
	if err := newTestClient(srv).get(context.Background(), "/repos/foo/bar", &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.FullName != want.FullName || got.Language != want.Language {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestGetDirContents_404_ReturnsNil(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	entries, err := newTestClient(srv).GetDirContents(context.Background(), "foo", "bar", ".github")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestGetCommits_StopsWhenPageEmpty(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls == 1 {
			json.NewEncoder(w).Encode([]Commit{{SHA: "a"}, {SHA: "b"}})
		} else {
			json.NewEncoder(w).Encode([]Commit{})
		}
	}))
	defer srv.Close()
	commits, err := newTestClient(srv).GetCommits(context.Background(), "foo", "bar", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 2 {
		t.Errorf("got %d commits, want 2", len(commits))
	}
}

func TestGetCommits_RespectsLimit(t *testing.T) {
	all := []Commit{{SHA: "a"}, {SHA: "b"}, {SHA: "c"}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perPage := len(all)
		if pp := r.URL.Query().Get("per_page"); pp != "" {
			fmt.Sscan(pp, &perPage)
		}
		page := all
		if perPage < len(page) {
			page = page[:perPage]
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}))
	defer srv.Close()
	commits, err := newTestClient(srv).GetCommits(context.Background(), "foo", "bar", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 2 {
		t.Errorf("got %d commits, want 2", len(commits))
	}
}
