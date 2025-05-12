package webhook

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/github"
)

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

	signature := r.Header.Get("X-Hub-Signature-256")
	if !VerifySignature(webhookSecret, body, signature) {
		http.Error(w, "Invald webhook signature", http.StatusForbidden)
		return
	}

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

	authenticator, err := github.NewAppAuthenticator(config.GithubAppId, config.GithubPrivateKeyPath)
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

	prNumber, err := strconv.ParseInt(config.PrNumber, 10, 64)
	if err != nil {
		http.Error(w, "Invalid PR number", http.StatusBadRequest)
		return
	}

	client := github.NewGitHubClient(token)
	err = client.PostComment(config.OwnerUsername, config.Repo, int(prNumber), "hello from your Github App!")
	if err != nil {
		log.Printf("Error posting comment: %v", err)
		http.Error(w, "Failed to post comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Comment posted successfully")
}
