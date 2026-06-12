package analyzer

import (
	"testing"

	"github.com/yannisobert/git-health-check/internal/github"
)

func minimalData() *github.RepoData {
	return &github.RepoData{
		Repo:         github.Repo{FullName: "foo/bar"},
		RootContents: []github.ContentEntry{{Name: "README.md", Size: 500}, {Name: "LICENSE"}},
		Commits:      []github.Commit{commit("feat: init", 5)},
	}
}

func TestBuildCategory_ScoreAggregation(t *testing.T) {
	checks := []CheckResult{
		{Score: 5, MaxScore: 10},
		{Score: 3, MaxScore: 5},
		{Score: 0, MaxScore: 5},
	}
	cat := buildCategory("Test", checks)
	if cat.Score != 8 {
		t.Errorf("Score = %d, want 8", cat.Score)
	}
	if cat.MaxScore != 20 {
		t.Errorf("MaxScore = %d, want 20", cat.MaxScore)
	}
	if cat.Name != "Test" {
		t.Errorf("Name = %q, want Test", cat.Name)
	}
}

func TestBuildReport_ScoreSumMatchesCategories(t *testing.T) {
	r := buildReport(minimalData(), "foo", "bar")

	var catSum, catMax int
	for _, c := range r.Categories {
		catSum += c.Score
		catMax += c.MaxScore
	}
	if r.Score != catSum {
		t.Errorf("Report.Score = %d, want %d (sum of categories)", r.Score, catSum)
	}
	if r.MaxScore != catMax {
		t.Errorf("Report.MaxScore = %d, want %d", r.MaxScore, catMax)
	}
}

func TestBuildReport_CategoryNames(t *testing.T) {
	r := buildReport(minimalData(), "foo", "bar")
	want := []string{"Maintenance", "Activity", "Conventions", "CI/CD"}
	if len(r.Categories) != len(want) {
		t.Fatalf("got %d categories, want %d", len(r.Categories), len(want))
	}
	for i, name := range want {
		if r.Categories[i].Name != name {
			t.Errorf("category[%d] = %q, want %q", i, r.Categories[i].Name, name)
		}
	}
}

func TestBuildReport_Metadata(t *testing.T) {
	data := minimalData()
	data.Repo.FullName = "owner/repo"
	r := buildReport(data, "owner", "repo")
	if r.FullName != "owner/repo" {
		t.Errorf("FullName = %q, want owner/repo", r.FullName)
	}
	if r.Owner != "owner" || r.Repo != "repo" {
		t.Errorf("Owner/Repo = %q/%q, want owner/repo", r.Owner, r.Repo)
	}
	if r.AnalyzedAt.IsZero() {
		t.Error("AnalyzedAt should not be zero")
	}
}
