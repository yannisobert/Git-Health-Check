package github

import "fmt"

var ErrRateLimit = fmt.Errorf("github rate limit exceeded")

type ErrAPI struct {
	Status int
}

func (e ErrAPI) Error() string {
	return fmt.Sprintf("github API error: status %d", e.Status)
}
