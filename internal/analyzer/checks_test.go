package analyzer

import (
	"testing"
	"time"

	"github.com/yannisobert/git-health-check/internal/github"
)

func commit(msg string, daysAgo int) github.Commit {
	var c github.Commit
	c.Commit.Message = msg
	c.Commit.Author.Date = time.Now().AddDate(0, 0, -daysAgo)
	return c
}

func entry(name string, size int) github.ContentEntry {
	return github.ContentEntry{Name: name, Size: size}
}

func pr(merged bool) github.PullRequest {
	p := github.PullRequest{State: "closed"}
	if merged {
		now := time.Now()
		p.MergedAt = &now
	}
	return p
}

func issue(state string, isPR bool) github.Issue {
	i := github.Issue{State: state}
	if isPR {
		i.PullRequest = &struct{}{}
	}
	return i
}

// --- Maintenance ---

func TestCheckREADME(t *testing.T) {
	tests := []struct {
		name     string
		entries  []github.ContentEntry
		wantScore int
	}{
		{"missing", nil, 0},
		{"short", []github.ContentEntry{entry("README.md", 50)}, 5},
		{"ok", []github.ContentEntry{entry("README.md", 500)}, 10},
		{"case insensitive", []github.ContentEntry{entry("readme.md", 500)}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkREADME(&github.RepoData{RootContents: tt.entries})
			if got.Score != tt.wantScore {
				t.Errorf("score = %d, want %d", got.Score, tt.wantScore)
			}
		})
	}
}

func TestCheckLICENSE(t *testing.T) {
	got := checkLICENSE(&github.RepoData{RootContents: []github.ContentEntry{entry("LICENSE", 0)}})
	if got.Score != 10 {
		t.Errorf("score = %d, want 10", got.Score)
	}
	got = checkLICENSE(&github.RepoData{})
	if got.Score != 0 {
		t.Errorf("score = %d, want 0", got.Score)
	}
}

func TestCheckCHANGELOG(t *testing.T) {
	got := checkCHANGELOG(&github.RepoData{RootContents: []github.ContentEntry{entry("CHANGELOG.md", 0)}})
	if got.Score != 5 {
		t.Errorf("score = %d, want 5", got.Score)
	}
}

// --- Activity ---

func TestCheckLastCommit(t *testing.T) {
	tests := []struct {
		name      string
		daysAgo   int
		wantScore int
	}{
		{"recent", 5, 10},
		{"month", 60, 7},
		{"quarter", 120, 4},
		{"old", 200, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &github.RepoData{Commits: []github.Commit{commit("feat: x", tt.daysAgo)}}
			got := checkLastCommit(data)
			if got.Score != tt.wantScore {
				t.Errorf("score = %d, want %d", got.Score, tt.wantScore)
			}
		})
	}
}

func TestCheckIssueRatio(t *testing.T) {
	data := &github.RepoData{Issues: []github.Issue{
		issue("closed", false),
		issue("closed", false),
		issue("closed", false),
		issue("open", false),
		issue("open", true), // PR, should be ignored
	}}
	got := checkIssueRatio(data)
	// 3 closed / 4 total = 75% → score 5
	if got.Score != 5 {
		t.Errorf("score = %d, want 5", got.Score)
	}
}

func TestCheckPRMergeRatio(t *testing.T) {
	data := &github.RepoData{PullRequests: []github.PullRequest{
		pr(true), pr(true), pr(true), pr(false),
	}}
	got := checkPRMergeRatio(data)
	// 3/4 = 75% → score 7
	if got.Score != 7 {
		t.Errorf("score = %d, want 7", got.Score)
	}
}

// --- Conventions ---

func TestCheckConventionalCommits(t *testing.T) {
	data := &github.RepoData{Commits: []github.Commit{
		commit("feat: add feature", 1),
		commit("fix: bug", 2),
		commit("random commit", 3),
		commit("chore: update deps", 4),
	}}
	got := checkConventionalCommits(data)
	// 3/4 = 75% → score 6
	if got.Score != 6 {
		t.Errorf("score = %d, want 6", got.Score)
	}
}

func TestCheckSemverReleases(t *testing.T) {
	data := &github.RepoData{Releases: []github.Release{
		{TagName: "v1.0.0"},
		{TagName: "v1.1.0"},
		{TagName: "not-semver"},
	}}
	got := checkSemverReleases(data)
	if got.Score != 5 {
		t.Errorf("score = %d, want 5", got.Score)
	}
}

func TestCheckLinter(t *testing.T) {
	got := checkLinter(&github.RepoData{RootContents: []github.ContentEntry{entry(".eslintrc.json", 0)}})
	if got.Score != 5 {
		t.Errorf("score = %d, want 5", got.Score)
	}
	got = checkLinter(&github.RepoData{})
	if got.Score != 0 {
		t.Errorf("score = %d, want 0", got.Score)
	}
}

// --- CI/CD ---

func TestCheckGitHubActions(t *testing.T) {
	got := checkGitHubActions(&github.RepoData{Workflows: []github.ContentEntry{entry("ci.yml", 0)}})
	if got.Score != 8 {
		t.Errorf("score = %d, want 8", got.Score)
	}
	got = checkGitHubActions(&github.RepoData{})
	if got.Score != 0 {
		t.Errorf("score = %d, want 0", got.Score)
	}
}

func TestCheckDependabot(t *testing.T) {
	got := checkDependabot(&github.RepoData{DotGithub: []github.ContentEntry{entry("dependabot.yml", 0)}})
	if got.Score != 4 {
		t.Errorf("score = %d, want 4", got.Score)
	}
}

func TestCheckCONTRIBUTING(t *testing.T) {
	tests := []struct {
		name      string
		entries   []github.ContentEntry
		wantScore int
	}{
		{"missing", nil, 0},
		{"short", []github.ContentEntry{entry("CONTRIBUTING.md", 50)}, 2},
		{"ok", []github.ContentEntry{entry("CONTRIBUTING.md", 500)}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkCONTRIBUTING(&github.RepoData{RootContents: tt.entries})
			if got.Score != tt.wantScore {
				t.Errorf("score = %d, want %d", got.Score, tt.wantScore)
			}
		})
	}
}

func TestCheckPRTemplate(t *testing.T) {
	got := checkPRTemplate(&github.RepoData{DotGithub: []github.ContentEntry{entry("pull_request_template.md", 0)}})
	if got.Score != 5 {
		t.Errorf("score = %d, want 5", got.Score)
	}
	got = checkPRTemplate(&github.RepoData{})
	if got.Score != 0 {
		t.Errorf("score = %d, want 0", got.Score)
	}
}

func TestCheckOtherCI(t *testing.T) {
	got := checkOtherCI(&github.RepoData{RootContents: []github.ContentEntry{entry(".travis.yml", 0)}})
	if got.Score != 4 {
		t.Errorf("score = %d, want 4", got.Score)
	}
	got = checkOtherCI(&github.RepoData{})
	if got.Score != 0 {
		t.Errorf("score = %d, want 0", got.Score)
	}
}

func TestCheckPreCommit(t *testing.T) {
	got := checkPreCommit(&github.RepoData{RootContents: []github.ContentEntry{entry(".pre-commit-config.yaml", 0)}})
	if got.Score != 4 {
		t.Errorf("score = %d, want 4", got.Score)
	}
	got = checkPreCommit(&github.RepoData{})
	if got.Score != 0 {
		t.Errorf("score = %d, want 0", got.Score)
	}
}
