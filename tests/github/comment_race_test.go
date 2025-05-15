package github_test

import (
	"sync"
	"testing"

	"github.com/yemiwebby/code-review-agent/internal/github/directwebhook"
)

func TestPostReviewCommentRace(t *testing.T) {
	t.Skip("Temporarily skipping test for this handler to avoid GitHub Post calls")

	var wg sync.WaitGroup
	concurrency := 10

	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()
			owner := "dummy"
			repo := "repo"
			pr := 1
			comment := "simualted AI comment"
			file := "file"
			line := 2
			oldPatch := "old Patch"

			directwebhook.PostReviewComment(owner, repo, pr, comment, file, line, oldPatch)
		}(i)
	}

	wg.Wait()
}
