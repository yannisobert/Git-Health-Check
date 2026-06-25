package analyzer

import "testing"

func reportWith(fullName string, checks []CheckResult) *Report {
	cat := buildCategory("Test", checks)
	return &Report{
		FullName:   fullName,
		Score:      cat.Score,
		MaxScore:   cat.MaxScore,
		Categories: []Category{cat},
	}
}

func TestBuildDiff_Repo1Wins(t *testing.T) {
	r1 := reportWith("foo/a", []CheckResult{{Name: "README", Score: 10, MaxScore: 10}})
	r2 := reportWith("foo/b", []CheckResult{{Name: "README", Score: 5, MaxScore: 10}})

	d := buildDiff(r1, r2)["README"]
	if d.Winner != "foo/a" {
		t.Errorf("Winner = %q, want foo/a", d.Winner)
	}
	if d.Delta != 5 {
		t.Errorf("Delta = %d, want 5", d.Delta)
	}
}

func TestBuildDiff_Repo2Wins(t *testing.T) {
	r1 := reportWith("foo/a", []CheckResult{{Name: "LICENSE", Score: 0, MaxScore: 10}})
	r2 := reportWith("foo/b", []CheckResult{{Name: "LICENSE", Score: 10, MaxScore: 10}})

	d := buildDiff(r1, r2)["LICENSE"]
	if d.Winner != "foo/b" {
		t.Errorf("Winner = %q, want foo/b", d.Winner)
	}
	if d.Delta != 10 {
		t.Errorf("Delta = %d, want 10", d.Delta)
	}
}

func TestBuildDiff_Tie(t *testing.T) {
	r1 := reportWith("foo/a", []CheckResult{{Name: "README", Score: 7, MaxScore: 10}})
	r2 := reportWith("foo/b", []CheckResult{{Name: "README", Score: 7, MaxScore: 10}})

	d := buildDiff(r1, r2)["README"]
	if d.Winner != "tie" {
		t.Errorf("Winner = %q, want tie", d.Winner)
	}
	if d.Delta != 0 {
		t.Errorf("Delta = %d, want 0", d.Delta)
	}
}

func TestBuildDiff_ScoresPreserved(t *testing.T) {
	r1 := reportWith("foo/a", []CheckResult{{Name: "README", Score: 8, MaxScore: 10}})
	r2 := reportWith("foo/b", []CheckResult{{Name: "README", Score: 3, MaxScore: 10}})

	d := buildDiff(r1, r2)["README"]
	if d.Repo1Score != 8 {
		t.Errorf("Repo1Score = %d, want 8", d.Repo1Score)
	}
	if d.Repo2Score != 3 {
		t.Errorf("Repo2Score = %d, want 3", d.Repo2Score)
	}
}

func TestBuildDiff_MissingCheckSkipped(t *testing.T) {
	r1 := reportWith("foo/a", []CheckResult{
		{Name: "README", Score: 10, MaxScore: 10},
		{Name: "LICENSE", Score: 5, MaxScore: 10},
	})
	r2 := reportWith("foo/b", []CheckResult{
		{Name: "README", Score: 5, MaxScore: 10},
	})

	diff := buildDiff(r1, r2)
	if _, ok := diff["LICENSE"]; ok {
		t.Error("LICENSE should not appear in diff when missing from r2")
	}
	if len(diff) != 1 {
		t.Errorf("diff len = %d, want 1", len(diff))
	}
}
