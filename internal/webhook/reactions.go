package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/yemiwebby/code-review-agent/internal/github"
	"github.com/yemiwebby/code-review-agent/internal/openai"
)

type ReactionsResponse struct {
	CommentID int
	Comment   *github.AIComment
	Status    int
	Err       error
}

func CheckReactionsHandler(w http.ResponseWriter, r *http.Request) {
	repo := r.URL.Query().Get("repo")

	github.Mu.Lock()
	if len(github.AIComments) == 0 {
		github.Mu.Unlock()
		http.Error(w, "No AI review comments found. Please ensure the AI agent has reviewed your PR before merging.", http.StatusPreconditionFailed)
		return
	}
	github.Mu.Unlock()

	commentIDStr := r.URL.Query().Get("comment_id")

	res := ValidateReactionsParams(repo, commentIDStr)
	if res.Err != nil {
		http.Error(w, res.Err.Error(), res.Status)
		return
	}

	up, down, err := github.FetchReactions(repo, res.CommentID)
	if err != nil {
		http.Error(w, "Failed to fetch reactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if up > 0 || down > 0 {
		fmt.Fprintln(w, "Reaction detected.")
	} else {
		http.Error(w, "No reaction found", http.StatusPreconditionFailed)
	}
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

	github.Mu.Lock()
	comment, ok := github.AIComments[commentID]
	github.Mu.Unlock()

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
