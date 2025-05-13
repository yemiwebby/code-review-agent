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
)

type PullRequestPayload struct {
	Action      string `json:"action"`
	PullRequest struct {
		Number int `json:"number"`
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"`
		Owner    struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}

func GithubAppHandler(w http.ResponseWriter, r *http.Request) {
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
	repo := payload.Repository.FullName
	prNumber := payload.PullRequest.Number

	// Log the extracted PR info
	log.Printf("Processing PR #%d in %s by %s", prNumber, repo, owner)

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

	// Post a dynamic comment to the PR
	client := github.NewGitHubClient(token)
	comment := fmt.Sprintf("Hello from your GitHub App! (PR #%d in %s)", prNumber, repo)
	err = client.PostComment(owner, repo, prNumber, comment)
	if err != nil {
		log.Printf("Error posting comment: %v", err)
		http.Error(w, "Failed to post comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Comment posted successfully")
}
