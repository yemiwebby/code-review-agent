package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/github"
	githubapp "github.com/yemiwebby/code-review-agent/internal/github/app"
	"github.com/yemiwebby/code-review-agent/internal/reviewer"
	githubhelper "github.com/yemiwebby/code-review-agent/internal/utils/githubHelper"
)

// type PullRequestPayload struct {
// 	Action      string `json:"action"`
// 	PullRequest struct {
// 		Number int `json:"number"`
// 	} `json:"pull_request"`
// 	Repository struct {
// 		FullName string `json:"full_name"`
// 		Owner    struct {
// 			Login string `json:"login"`
// 		} `json:"owner"`
// 	} `json:"repository"`
// }

func GithubAppHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received GitHub App webhook")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	webhookSecret := config.GithubWebhookSecret
	if webhookSecret == "" {
		http.Error(w, "Webhook secret not set", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify the webhook signature
	signature := r.Header.Get("X-Hub-Signature-256")
	if !VerifySignature(webhookSecret, body, signature) {
		http.Error(w, "Invalid webhook signature", http.StatusForbidden)
		return
	}

	// Parse the installation ID
	installationIdStr := r.Header.Get("X-Github-Installation-Id")
	if installationIdStr == "" {
		http.Error(w, "Missing Installation ID", http.StatusBadRequest)
		return
	}

	installationID, err := strconv.ParseInt(installationIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid installation ID", http.StatusBadRequest)
		return
	}

	// Parse the webhook payload
	var payload PullRequestPayload
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		http.Error(w, "Failed to parse JSON payload", http.StatusBadRequest)
		return
	}

	owner := payload.Repository.Owner.Login
	repo, err := githubhelper.ExtractRepoName(owner, payload.Repository.FullName)
	if err != nil {
		log.Fatalf("Failed to extract repo name: %v", err)
	}

	prNumber := payload.PullRequest.Number

	// Log the extracted PR info
	log.Printf("Extracted PR info: owner=%s, extractedRepo=%s, prNumber=%d", owner, repo, prNumber)

	// Authenticate with GitHub
	authenticator, err := github.NewAppAuthenticator(config.GithubAppId)
	if err != nil {
		log.Printf("Error getting installation token: %v", err)
		http.Error(w, "Failed to initialize authenticator", http.StatusInternalServerError)
		return
	}

	token, err := authenticator.GetInstallationToken(installationID)
	if err != nil {
		log.Printf("Error getting installation token: %v", err)
		http.Error(w, "Failed to get installation token", http.StatusInternalServerError)
		return
	}

	// Review the PR using AI
	log.Printf("Starting AI review for PR #%d in %s/%s", prNumber, owner, repo)
	client := githubapp.NewGitHubAppClient(token)

	commitID, err := github.GetCommitID(owner, repo, payload.PullRequest.Number, token, false)
	if err != nil {
		log.Printf("failed to get commit ID: %v", err)
		return
	}

	go reviewer.ReviewPullRequest(owner, repo, prNumber, commitID, client)

	log.Println("AI review completed")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "AI review completed")
}
