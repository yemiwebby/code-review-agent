package reviewer

import (
	"fmt"
	"strings"

	"github.com/yemiwebby/code-review-agent/internal/github"
	"github.com/yemiwebby/code-review-agent/internal/github/directwebhook"
	"github.com/yemiwebby/code-review-agent/internal/openai"
)

func ReviewPullRequest(owner, repo string, prNumber int) {
	fmt.Printf("Reviewing PR #%d in %s/%s\n", prNumber, owner, repo)

	files, err := github.GetPRFiles(owner, repo, prNumber)
	if err != nil {
		fmt.Println("Failed to fetch PR files:", err)
		return
	}

	analyzeAndPostComments(files, owner, repo, prNumber)
	processCommentReactions(owner, repo, prNumber)
}

func analyzeAndPostComments(files []github.FileChange, owner, repo string, prNumber int) {
	for _, file := range files {
		if file.Patch == "" {
			continue
		}

		if !strings.HasSuffix(file.Filename, ".go") {
			fmt.Printf("Skipping non-Go file: %s\n", file.Filename)
			continue
		}

		review, err := openai.AnalyzeCode(file.Patch)
		if err != nil {
			fmt.Println("AI Review failed:", err)
			continue
		}

		comment := fmt.Sprintf("**File:** %s:\n%s\n\n", file.Filename, review)

		err = directwebhook.PostReviewComment(owner, repo, prNumber, comment, file.Filename, 0, file.Patch)
		if err != nil {
			fmt.Println("Failed to post comment:", err)
		}
	}
}

func processCommentReactions(owner, repo string, prNumber int) {
	directwebhook.Mu.Lock()
	defer directwebhook.Mu.Unlock()

	for id, aiComment := range directwebhook.AIComments {
		up, down, err := github.FetchReactions(repo, id)
		if err != nil {
			fmt.Println("Could not fetch reactions:", err)
			continue
		}

		if up == 0 && down == 0 {
			fmt.Println("No reactions yet, skipping prompt adjustment.")
			continue
		}

		if aiComment.OldPatch == "" {
			fmt.Printf("No patch found for comment %d\n", id)
			continue
		}

		changed := hasCodeChanged(owner, repo, prNumber, aiComment.File, aiComment.OldPatch)

		if !changed {
			fmt.Printf("Code hasn't changed since AI comment for file: %s\n", aiComment.File)
		} else {
			fmt.Printf("Code changed for: %s\n", aiComment.File)
		}

		adjusted := openai.AdjustPrompt(aiComment.Body, up, down)
		fmt.Println("Adjusted prompt:\n", adjusted)
	}

}

func hasCodeChanged(owner, repo string, prNumber int, filename, oldPatch string) bool {
	currentFiles, err := github.GetPRFiles(owner, repo, prNumber)
	if err != nil {
		fmt.Println("Failed to fetch latest PR files:", err)
		return false
	}

	for _, f := range currentFiles {
		if f.Filename == filename && f.Patch != oldPatch {
			return true
		}
	}
	return false
}
