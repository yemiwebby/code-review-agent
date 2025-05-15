package github_test

import (
	"sync"
	"testing"

	"github.com/yemiwebby/code-review-agent/internal/github"
	"github.com/yemiwebby/code-review-agent/internal/github/directwebhook"
)

func TestPostReviewCommentRace(t *testing.T) {
	t.Skip("Temporarily skipping test for this handler to avoid GitHub Post calls")

	var wg sync.WaitGroup
	concurrency := 10

	// Use the CommentPoster interface for flexibility
	var client github.GithubClientInterface = &directwebhook.GithubClient{}

	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()
			owner := "dummy"
			repo := "repo"
			pr := 1
			comment := "simulated AI comment"
			file := "file.go"
			line := 2
			oldPatch := "old Patch"

			err := client.PostReviewComment(owner, repo, pr, comment, file, "", line, oldPatch)
			if err != nil {
				t.Errorf("Failed to post review comment: %v", err)
			}
		}(i)
	}

	wg.Wait()
}
