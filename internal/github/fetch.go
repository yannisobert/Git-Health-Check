package github

import (
	"context"
	"sync"
)

// FetchAll retrieves all repository data in parallel.
func (c *Client) FetchAll(ctx context.Context, owner, repo string) (*RepoData, error) {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		data RepoData
		errs []error
	)

	fail := func(err error) {
		mu.Lock()
		errs = append(errs, err)
		mu.Unlock()
	}

	set := func(fn func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn()
		}()
	}

	set(func() {
		r, err := c.GetRepo(ctx, owner, repo)
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.Repo = r
		mu.Unlock()
	})

	set(func() {
		commits, err := c.GetCommits(ctx, owner, repo, 30)
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.Commits = commits
		mu.Unlock()
	})

	set(func() {
		issues, err := c.GetIssues(ctx, owner, repo)
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.Issues = issues
		mu.Unlock()
	})

	set(func() {
		prs, err := c.GetPullRequests(ctx, owner, repo)
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.PullRequests = prs
		mu.Unlock()
	})

	set(func() {
		releases, err := c.GetReleases(ctx, owner, repo)
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.Releases = releases
		mu.Unlock()
	})

	set(func() {
		langs, err := c.GetLanguages(ctx, owner, repo)
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.Languages = langs
		mu.Unlock()
	})

	set(func() {
		entries, err := c.GetDirContents(ctx, owner, repo, "")
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.RootContents = entries
		mu.Unlock()
	})

	set(func() {
		entries, err := c.GetDirContents(ctx, owner, repo, ".github")
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.DotGithub = entries
		mu.Unlock()
	})

	set(func() {
		entries, err := c.GetDirContents(ctx, owner, repo, ".github/workflows")
		if err != nil {
			fail(err)
			return
		}
		mu.Lock()
		data.Workflows = entries
		mu.Unlock()
	})

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}
	return &data, nil
}
