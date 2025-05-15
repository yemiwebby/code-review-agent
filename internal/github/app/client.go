package githubapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yemiwebby/code-review-agent/internal/comment"
)

type GitHubAppClient struct {
	Token string
}

func NewGitHubAppClient(token string) *GitHubAppClient {
	return &GitHubAppClient{Token: token}
}

func (c *GitHubAppClient) PostReviewComment(owner, repo string, prNumber int, body, file, commitID string, line int, patch string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, prNumber)

	payload := map[string]interface{}{
		"body":      body,
		"commit_id": commitID,
		"path":      file,
		"positon":   line,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal comment data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post comment, status code: %d, response: %s", resp.StatusCode, respBody)
	}

	// Extract the comment ID for reaction tracking
	var result struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to parse comment ID from response: %w", err)
	}

	// Store the comment for future reference
	commentID := result.ID
	comment.StoreComment(commentID, body, file, line, patch)

	fmt.Printf("Posted review comment for %s: %s (ID: %d)\n", file, body, commentID)
	return nil
}
