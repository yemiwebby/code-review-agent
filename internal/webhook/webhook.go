package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yemiwebby/code-review-agent/internal/reviewer"
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
}

type PullRequestEvent struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var event PullRequestEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if event.Action == "opened" || event.Action == "synchronize" {
		go reviewer.ReviewPullRequest(event.Repository.FullName, event.PullRequest.Number)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "Received webhook")
}
