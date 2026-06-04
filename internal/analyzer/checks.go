package analyzer

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yannisobert/git-health-check/internal/github"
)

var semverRe = regexp.MustCompile(`^v\d+\.\d+\.\d+`)

var conventionalPrefixes = []string{
	"feat:", "fix:", "chore:", "docs:", "test:",
	"refactor:", "style:", "perf:", "ci:", "build:", "revert:",
}

var linterFiles = []string{
	".eslintrc", ".eslintrc.json", ".eslintrc.js", ".eslintrc.cjs", ".eslintrc.yml", ".eslintrc.yaml",
	"eslint.config.js", "eslint.config.mjs", "eslint.config.cjs",
	".prettierrc", ".prettierrc.json", ".prettierrc.yaml", ".prettierrc.yml", ".prettierrc.js",
	"golangci.yml", ".golangci.yml", ".golangci.yaml",
	".rubocop.yml",
	"pyproject.toml",
	".flake8",
	"biome.json",
}

var otherCIFiles = []string{
	".travis.yml", "jenkinsfile", "azure-pipelines.yml", ".circleci", "circle.yml",
}

var preCommitFiles = []string{
	".husky", ".pre-commit-config.yaml", ".pre-commit-config.yml", ".lefthook.yml",
}

// Helpers

func findFile(entries []github.ContentEntry, names ...string) (github.ContentEntry, bool) {
	for _, e := range entries {
		lower := strings.ToLower(e.Name)
		for _, n := range names {
			if lower == n {
				return e, true
			}
		}
	}
	return github.ContentEntry{}, false
}

func ok(name string, max int, detail string) CheckResult {
	return CheckResult{Name: name, Score: max, MaxScore: max, Status: StatusOK, Detail: detail}
}

func fail(name string, max int, detail, suggestion string) CheckResult {
	return CheckResult{Name: name, Score: 0, MaxScore: max, Status: StatusFail, Detail: detail, Suggestion: suggestion}
}

func partial(name string, score, max int, detail, suggestion string) CheckResult {
	st := StatusWarn
	if score == 0 {
		st = StatusFail
	}
	return CheckResult{Name: name, Score: score, MaxScore: max, Status: st, Detail: detail, Suggestion: suggestion}
}

// Maintenance (30 pts)

func checkREADME(data *github.RepoData) CheckResult {
	e, found := findFile(data.RootContents, "readme.md", "readme", "readme.txt", "readme.rst")
	if !found {
		return fail("README", 10, "No README found", "Add a README.md describing the project")
	}
	if e.Size < 200 {
		return partial("README", 5, 10,
			fmt.Sprintf("README found but short (%d bytes)", e.Size),
			"Expand your README with installation, usage, and contribution instructions")
	}
	return ok("README", 10, fmt.Sprintf("README found (%d bytes)", e.Size))
}

func checkLICENSE(data *github.RepoData) CheckResult {
	if _, found := findFile(data.RootContents, "license", "license.md", "license.txt", "licence", "copying"); found {
		return ok("LICENSE", 10, "LICENSE file found")
	}
	return fail("LICENSE", 10, "No LICENSE file found", "Add a LICENSE file to clarify how the project can be used")
}

func checkCHANGELOG(data *github.RepoData) CheckResult {
	if _, found := findFile(data.RootContents, "changelog.md", "changelog", "changelog.txt", "history.md"); found {
		return ok("CHANGELOG", 5, "CHANGELOG found")
	}
	return fail("CHANGELOG", 5, "No CHANGELOG found", "Add a CHANGELOG.md to document version history")
}

func checkCONTRIBUTING(data *github.RepoData) CheckResult {
	e, found := findFile(data.RootContents, "contributing.md", "contributing", "contributing.txt")
	if !found {
		return fail("CONTRIBUTING", 5, "No CONTRIBUTING guide found", "Add a CONTRIBUTING.md to help new contributors")
	}
	if e.Size < 100 {
		return partial("CONTRIBUTING", 2, 5,
			"CONTRIBUTING found but very short",
			"Expand your CONTRIBUTING guide with development setup and PR process")
	}
	return ok("CONTRIBUTING", 5, fmt.Sprintf("CONTRIBUTING found (%d bytes)", e.Size))
}

// --- Activity (25 pts) ---

func checkLastCommit(data *github.RepoData) CheckResult {
	if len(data.Commits) == 0 {
		return fail("Last commit", 10, "No commits found", "")
	}
	days := int(time.Since(data.Commits[0].Commit.Author.Date).Hours() / 24)
	detail := fmt.Sprintf("Last commit %d days ago", days)
	switch {
	case days < 30:
		return ok("Last commit", 10, detail)
	case days < 90:
		return partial("Last commit", 7, 10, detail, "Aim for regular commits to show active maintenance")
	case days < 180:
		return partial("Last commit", 4, 10, detail, "Repository appears inactive — consider updating or archiving")
	default:
		return fail("Last commit", 10, detail, "Repository appears abandoned — consider archiving it")
	}
}

func checkIssueRatio(data *github.RepoData) CheckResult {
	var open, closed int
	for _, i := range data.Issues {
		if i.IsPR() {
			continue
		}
		if i.State == "closed" {
			closed++
		} else {
			open++
		}
	}
	total := open + closed
	if total == 0 {
		return partial("Issue closure rate", 4, 8, "No issues found", "")
	}
	ratio := float64(closed) / float64(total)
	detail := fmt.Sprintf("%d/%d issues closed (%.0f%%)", closed, total, ratio*100)
	switch {
	case ratio >= 0.8:
		return ok("Issue closure rate", 8, detail)
	case ratio >= 0.5:
		return partial("Issue closure rate", 5, 8, detail, "Work on closing open issues to improve project health")
	case ratio >= 0.2:
		return partial("Issue closure rate", 2, 8, detail, "Many open issues — consider triaging and closing stale ones")
	default:
		return fail("Issue closure rate", 8, detail, "Most issues are unaddressed — triage and close stale issues")
	}
}

func checkPRMergeRatio(data *github.RepoData) CheckResult {
	total := len(data.PullRequests)
	if total == 0 {
		return partial("PR merge rate", 3, 7, "No pull requests found", "")
	}
	var merged int
	for _, pr := range data.PullRequests {
		if pr.MergedAt != nil {
			merged++
		}
	}
	ratio := float64(merged) / float64(total)
	detail := fmt.Sprintf("%d/%d PRs merged (%.0f%%)", merged, total, ratio*100)
	switch {
	case ratio >= 0.7:
		return ok("PR merge rate", 7, detail)
	case ratio >= 0.4:
		return partial("PR merge rate", 4, 7, detail, "Consider reviewing and merging or closing open PRs")
	default:
		return fail("PR merge rate", 7, detail, "Low PR merge rate — triage open pull requests")
	}
}

// Conventions (25 pts)

func checkConventionalCommits(data *github.RepoData) CheckResult {
	if len(data.Commits) == 0 {
		return fail("Conventional commits", 10, "No commits to analyze", "")
	}
	var matching int
	for _, c := range data.Commits {
		msg := strings.ToLower(c.Commit.Message)
		for _, prefix := range conventionalPrefixes {
			if strings.HasPrefix(msg, prefix) {
				matching++
				break
			}
		}
	}
	ratio := float64(matching) / float64(len(data.Commits))
	detail := fmt.Sprintf("%d/%d commits follow conventional format (%.0f%%)", matching, len(data.Commits), ratio*100)
	switch {
	case ratio >= 0.8:
		return ok("Conventional commits", 10, detail)
	case ratio >= 0.5:
		return partial("Conventional commits", 6, 10, detail, "Adopt conventional commits format consistently (feat:, fix:, chore:...)")
	case ratio >= 0.2:
		return partial("Conventional commits", 3, 10, detail, "Use conventional commits format consistently")
	default:
		return fail("Conventional commits", 10, detail, "Adopt conventional commits: feat:, fix:, chore:, docs:, test:...")
	}
}

func checkPRTemplate(data *github.RepoData) CheckResult {
	if _, found := findFile(data.DotGithub, "pull_request_template.md"); found {
		return ok("PR template", 5, "Pull request template found")
	}
	return fail("PR template", 5, "No PR template found", "Add .github/PULL_REQUEST_TEMPLATE.md to standardize pull requests")
}

func checkSemverReleases(data *github.RepoData) CheckResult {
	if len(data.Releases) == 0 {
		return fail("Semver releases", 5, "No releases found", "Publish versioned releases tagged as vX.Y.Z")
	}
	var count int
	for _, r := range data.Releases {
		if semverRe.MatchString(r.TagName) {
			count++
		}
	}
	if count == 0 {
		return fail("Semver releases", 5,
			fmt.Sprintf("%d release(s) found but none follow semver", len(data.Releases)),
			"Tag releases as vX.Y.Z (e.g. v1.2.3)")
	}
	return ok("Semver releases", 5, fmt.Sprintf("%d semver release(s) found", count))
}

func checkLinter(data *github.RepoData) CheckResult {
	for _, e := range data.RootContents {
		lower := strings.ToLower(e.Name)
		for _, f := range linterFiles {
			if lower == f {
				return ok("Linter/formatter", 5, fmt.Sprintf("Linter config found: %s", e.Name))
			}
		}
	}
	return fail("Linter/formatter", 5, "No linter or formatter config found", "Add a linter config (.eslintrc, golangci.yml, .rubocop.yml...)")
}

// CI/CD (20 pts)

func checkGitHubActions(data *github.RepoData) CheckResult {
	if len(data.Workflows) > 0 {
		return ok("GitHub Actions", 8, fmt.Sprintf("%d workflow file(s) found", len(data.Workflows)))
	}
	return fail("GitHub Actions", 8, "No GitHub Actions workflows found", "Add CI workflows in .github/workflows/")
}

func checkOtherCI(data *github.RepoData) CheckResult {
	for _, e := range data.RootContents {
		lower := strings.ToLower(e.Name)
		for _, f := range otherCIFiles {
			if lower == f {
				return ok("Other CI", 4, fmt.Sprintf("CI config found: %s", e.Name))
			}
		}
	}
	return fail("Other CI", 4, "No external CI config found", "Consider adding CI (Travis, CircleCI, Jenkins...)")
}

func checkPreCommit(data *github.RepoData) CheckResult {
	for _, e := range data.RootContents {
		lower := strings.ToLower(e.Name)
		for _, f := range preCommitFiles {
			if lower == f {
				return ok("Pre-commit hooks", 4, fmt.Sprintf("Pre-commit config found: %s", e.Name))
			}
		}
	}
	return fail("Pre-commit hooks", 4, "No pre-commit hooks found", "Add .pre-commit-config.yaml or .husky for automated checks")
}

func checkDependabot(data *github.RepoData) CheckResult {
	if _, found := findFile(data.DotGithub, "dependabot.yml", "dependabot.yaml"); found {
		return ok("Dependabot", 4, "Dependabot config found")
	}
	return fail("Dependabot", 4, "No Dependabot config found", "Add .github/dependabot.yml to automate dependency updates")
}
