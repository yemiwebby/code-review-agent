package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/github"
	"github.com/yemiwebby/code-review-agent/internal/github/directwebhook"
	"github.com/yemiwebby/code-review-agent/internal/reviewer"
	githubhelper "github.com/yemiwebby/code-review-agent/internal/utils/githubHelper"
)

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	User   User   `json:"user"`
}

type User struct {
	Login string `json:"login"`
}

type Repository struct {
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
}

type PullRequestPayload struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload PullRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// if event.Action == "opened" || event.Action == "synchronize" {
	owner, repo, err := githubhelper.SplitRepoName(payload.Repository.FullName)
	if err != nil {
		log.Fatalf("Failed to extract repo name: %v", err)
	}

	token := config.GithubToken
	client := directwebhook.NewGitHubClient(token)
	commitID, err := github.GetCommitID(owner, repo, payload.PullRequest.Number, token, false)
	if err != nil {
		log.Printf("failed to get commit ID: %v", err)
		return
	}

	go reviewer.ReviewPullRequest(owner, repo, payload.PullRequest.Number, commitID, client)
	// }

	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "Received webhook")
}
