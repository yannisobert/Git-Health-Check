package analyzer

import (
	"context"
	"testing"
)

func TestRivals_KnownRepo(t *testing.T) {
	a := &Analyzer{client: nil} // client never called for known repos
	suggestions, err := a.Rivals(context.Background(), "gin-gonic", "gin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(suggestions) == 0 {
		t.Fatal("expected suggestions for known repo")
	}
	for _, s := range suggestions {
		if s.FullName == "" {
			t.Error("expected non-empty FullName")
		}
		if s.Similarity != "direct-rival" {
			t.Errorf("Similarity = %q, want direct-rival", s.Similarity)
		}
	}
}

func TestRivals_KnownRepo_ReturnsCorrectRivals(t *testing.T) {
	a := &Analyzer{client: nil}
	suggestions, err := a.Rivals(context.Background(), "facebook", "react")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := make(map[string]bool)
	for _, s := range suggestions {
		names[s.FullName] = true
	}
	if !names["vuejs/vue"] {
		t.Error("expected vuejs/vue in react rivals")
	}
	if !names["sveltejs/svelte"] {
		t.Error("expected sveltejs/svelte in react rivals")
	}
}
