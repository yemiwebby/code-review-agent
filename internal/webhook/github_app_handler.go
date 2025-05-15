package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/github"
	githubapp "github.com/yemiwebby/code-review-agent/internal/github/app"
	"github.com/yemiwebby/code-review-agent/internal/openai"
	githubhelper "github.com/yemiwebby/code-review-agent/internal/utils/githubHelper"
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
	repo := payload.Repository.FullName

	extractedRepo, err := githubhelper.ExtractRepoName(owner, repo)
	if err != nil {
		log.Fatalf("Failed to extract repo name: %v", err)
	}

	prNumber := payload.PullRequest.Number

	// Log the extracted PR info
	log.Printf("Extracted PR info: owner=%s, extractedRepo=%s, prNumber=%d", owner, extractedRepo, prNumber)

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
	log.Printf("Strting AI review for PR #%d in %s", prNumber, repo)
	client := githubapp.NewGitHubAppClient(token)
	files, err := github.GetPRFiles(owner, extractedRepo, prNumber)
	if err != nil {
		log.Printf("Failed to fetch PR files: %v", err)
		http.Error(w, "Failed to fetch PR files", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if file.Patch == "" || !strings.HasSuffix(file.Filename, ".go") {
			log.Printf("Skipping non-Go file or empty patch: %s", file.Filename)
			continue
		}

		review, err := openai.AnalyzeCode(file.Patch)
		if err != nil {
			log.Printf("AI review failed for %s: %v", file.Filename, err)
			continue
		}

		comment := fmt.Sprintf("**File:** %s\n\n%s\n\n", file.Filename, review)
		log.Printf("Posting comment: %s", comment)

		err = client.PostReviewComment(owner, extractedRepo, prNumber, comment, file.Filename, 0, file.Patch)
		if err != nil {
			log.Printf("Failed to post comment for %s: %v", file.Filename, err)
		}
	}

	log.Println("AI review completed")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "AI review completed")
}
