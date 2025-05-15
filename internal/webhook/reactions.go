package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yemiwebby/code-review-agent/internal/comment"
	"github.com/yemiwebby/code-review-agent/internal/github"
	"github.com/yemiwebby/code-review-agent/internal/openai"
)

const reactionGracePeriod = 30 * time.Second

type ReactionsResponse struct {
	CommentID int
	Comment   *comment.AIComment
	Status    int
	Err       error
}

func CheckReactionsHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	if repo == "" {
		http.Error(w, "Missing 'repo' query parameter", http.StatusBadRequest)
		return
	}

	comment.Mu.Lock()
	defer comment.Mu.Unlock()

	if len(comment.AIComments) == 0 {
		http.Error(w, "No AI review comments found. Ensure AI agent has reviewed the PR.", http.StatusPreconditionFailed)
		return
	}

	var total int
	var enforced int
	var acknowledged int
	now := time.Now()

	for id, comment := range comment.AIComments {
		total++

		if now.Sub(comment.Timestamp) < reactionGracePeriod {
			fmt.Printf("Skipping recent comment %d (%s ago)\n", id, now.Sub(comment.Timestamp).Round(time.Second))
			continue
		}

		if comment.OldPatch == "" {
			fmt.Printf("Skipping comment %d due to missing OldPatch\n", id)
			continue
		}

		up, down, err := github.FetchReactions(repo, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch reactions for comment %d: %v", id, err), http.StatusInternalServerError)
			return
		}

		enforced++

		if up > 0 || down > 0 {
			acknowledged++
		} else {
			fmt.Printf("No reaction on comment %d (%s)\n", id, comment.File)
		}
	}

	if enforced == 0 {
		http.Error(w, "No eligible comments found for enforcement yet. Try again ;ater", http.StatusPreconditionFailed)
		return
	}

	if acknowledged == 0 {
		http.Error(w, fmt.Sprintf("No reactions found on %d enforced comment(s).", enforced), http.StatusPreconditionFailed)
	}

	if acknowledged == 0 {
		http.Error(w, fmt.Sprintf("No reactions found on %d enforced comment(s).", enforced), http.StatusPreconditionFailed)
		return
	}

	fmt.Fprintln(w, "Enforced AI comments have been acknowledged with a reaction.\n", acknowledged, enforced)
}

func ProcessReactions(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")
	idStr := r.URL.Query().Get("comment_id")

	res := ValidateReactionsParams(repo, idStr)
	if res.Err != nil {
		http.Error(w, res.Err.Error(), res.Status)
		return
	}

	up, down, err := github.FetchReactions(repo, res.CommentID)
	if err != nil {
		http.Error(w, "Failed to fetch reactions", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "👍 %d | 👎 %d for comment %d\n", up, down, res.CommentID)

	adjusted := openai.AdjustPrompt(res.Comment.Body, up, down)
	fmt.Fprintf(w, "\n Adjusted Prompt:\n%s\n", adjusted)
}

func ValidateReactionsParams(repo, commentIDStr string) ReactionsResponse {

	if repo == "" || commentIDStr == "" {
		return NewReactionsErr(http.StatusBadRequest, errors.New("missing 'repo' or 'comment_id' query parameter"))
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return NewReactionsErr(http.StatusBadRequest, errors.New("invalid comment_id"))
	}

	comment.Mu.Lock()
	defer comment.Mu.Unlock()

	comment, ok := comment.AIComments[commentID]
	if !ok {
		return NewReactionsErr(http.StatusNotFound, errors.New("comment not found in memory"))
	}

	if !ok {
		return NewReactionsErr(http.StatusNotFound, errors.New("comment not found in memory"))
	}

	return ReactionsResponse{
		CommentID: commentID,
		Comment:   comment,
		Status:    http.StatusOK,
		Err:       nil,
	}
}

func NewReactionsErr(status int, err error) ReactionsResponse {
	return ReactionsResponse{
		Status: status,
		Err:    err,
	}
}
