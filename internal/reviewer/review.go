package reviewer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/yemiwebby/code-review-agent/internal/comment"
	"github.com/yemiwebby/code-review-agent/internal/github"
	"github.com/yemiwebby/code-review-agent/internal/openai"
)

func ReviewPullRequest(owner, repo string, prNumber int, commitID string, client github.GithubClientInterface) {
	fmt.Printf("Reviewing PR #%d in %s/%s\n", prNumber, owner, repo)

	files, err := github.GetPRFiles(owner, repo, prNumber)
	if err != nil {
		fmt.Printf("Failed to fetch PR files for %s/%s: %s\n", owner, repo, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, file := range files {
		go func(file github.FileChange) {
			defer wg.Done()
			processFileChange(file, owner, repo, prNumber, commitID, client)
		}(file)
	}

	wg.Wait()

	// analyzeAndPostComments(files, owner, repo, prNumber, commitID, client)
	processCommentReactions(owner, repo, prNumber, client)

	fmt.Println("Review completed for all files")
}

func processFileChange(file github.FileChange, owner, repo string, prNumber int, commitID string, client github.GithubClientInterface) {
	if file.Patch == "" {
		fmt.Printf("Skipping file without patch: %s\n", file.Filename)
		return
	}

	if !strings.HasSuffix(file.Filename, ".go") {
		fmt.Printf("Skipping non-Go file: %s\n", file.Filename)
		return
	}

	review, err := openai.AnalyzeCode(file.Patch)
	if err != nil {
		fmt.Printf("AI Review failed for %s: %s\n", file.Filename, err)
		return
	}

	position, err := extractLinePosition(file.Patch)
	if err != nil {
		fmt.Printf("Failed to extract line positio for %s: %s\n", file.Filename, err)
		return
	}

	commentBody := fmt.Sprintf("**File:** %s:\n%s\n\n", file.Filename, review)

	commentID, err := client.PostReviewComment(owner, repo, prNumber, commentBody, file.Filename, commitID, position, file.Patch)
	if err != nil {
		fmt.Printf("Failed to post comment for %s: %s\n", file.Filename, err)
		return
	}

	comment.StoreComment(commentID, commentBody, file.Filename, position, file.Patch)
	fmt.Printf("Posted review comment for %s (commit ID: %s)\n", file.Filename, commitID)
}

func processCommentReactions(owner, repo string, prNumber int, client github.GithubClientInterface) {
	comment.Mu.Lock()
	defer comment.Mu.Unlock()

	for id, aiComment := range comment.AIComments {

		if aiComment.File == "" || aiComment.OldPatch == "" {
			fmt.Printf("Skipping comment %d due to missing metadata (file: %s, patch: %s)\n", id, aiComment.File, aiComment.OldPatch)
			continue
		}

		up, down, err := client.FetchReactions(owner, repo, id)
		if err != nil {
			fmt.Printf("Could not fetch reactions for comment %d: %s\n", id, err)
			continue
		}

		if up == 0 && down == 0 {
			fmt.Printf("No reactions yet for comment %d, skipping prompt adjustment.\n", id)
			continue
		}

		if aiComment.OldPatch == "" {
			fmt.Printf("No patch found for comment %d, skipping code change check.\n", id)
			continue
		}

		if !hasCodeChanged(owner, repo, prNumber, aiComment.File, aiComment.OldPatch) {
			fmt.Printf("Code hasn't changed since AI comment for file: %s (ID: %d)\n", aiComment.File, id)
			continue
		}

		fmt.Printf("Code changed for: %s (ID: %d)\n", aiComment.File, id)

		adjusted := openai.AdjustPrompt(aiComment.Body, up, down)
		fmt.Printf("Adjusted prompt for comment %d:\n%s\n", id, adjusted)
	}

}

func hasCodeChanged(owner, repo string, prNumber int, filename, oldPatch string) bool {
	currentFiles, err := github.GetPRFiles(owner, repo, prNumber)
	if err != nil {
		fmt.Printf("Failed to fetch latest PR files for %s/%s: %s\n", owner, repo, err)
		return false
	}

	for _, f := range currentFiles {
		if f.Filename == filename && f.Patch != oldPatch {
			return true
		}
	}
	return false
}

func extractLinePosition(patch string) (int, error) {
	lines := strings.Split(patch, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++") {
			return i + 1, nil
		}
	}

	return 1, nil
}
