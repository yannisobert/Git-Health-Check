package analyzer

import (
	"context"
	"sync"
)

func (a *Analyzer) Compare(ctx context.Context, owner1, repo1, owner2, repo2 string) (*CompareResult, error) {
	var (
		r1, r2     *Report
		err1, err2 error
		wg         sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		r1, err1 = a.Analyze(ctx, owner1, repo1)
	}()
	go func() {
		defer wg.Done()
		r2, err2 = a.Analyze(ctx, owner2, repo2)
	}()
	wg.Wait()

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	return &CompareResult{
		Repo1: *r1,
		Repo2: *r2,
		Diff:  buildDiff(r1, r2),
	}, nil
}

func buildDiff(r1, r2 *Report) map[string]CheckDiff {
	r2idx := make(map[string]CheckResult)
	for _, cat := range r2.Categories {
		for _, c := range cat.Checks {
			r2idx[c.Name] = c
		}
	}

	diff := make(map[string]CheckDiff)
	for _, cat := range r1.Categories {
		for _, c := range cat.Checks {
			c2, ok := r2idx[c.Name]
			if !ok {
				continue
			}
			delta := c.Score - c2.Score
			winner := "tie"
			if delta > 0 {
				winner = r1.FullName
			} else if delta < 0 {
				winner = r2.FullName
				delta = -delta
			}
			diff[c.Name] = CheckDiff{
				Repo1Score: c.Score,
				Repo2Score: c2.Score,
				Delta:      delta,
				Winner:     winner,
			}
		}
	}
	return diff
}
