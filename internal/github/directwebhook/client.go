package directwebhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yemiwebby/code-review-agent/internal/github"
)

type GithubClient struct {
	Token string
}

func NewGitHubClient(token string) *GithubClient {
	return &GithubClient{Token: token}
}

func (c *GithubClient) PostReviewComment(owner, repo string, prNumber int, body, file, commitID string, line int, patch string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, prNumber)

	payload := map[string]interface{}{
		"body":      body,
		"commit_id": commitID,
		"path":      file,
		"position":  line,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal comment data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("failed to post comment, status code: %d, response: %s", resp.StatusCode, respBody)
	}

	var result struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return 0, fmt.Errorf("failed to parse comment ID from response: %w", err)
	}

	fmt.Printf("Posted review comment for %s: %s (ID: %d)\n", file, body, result.ID)
	return result.ID, nil
}

func (c *GithubClient) FetchReactions(owner, repo string, commentID int) (int, int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/comments/%d/reactions", owner, repo, commentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "token "+c.Token)
	req.Header.Set("Accept", "application/vnd.github.squirrel-girl-preview+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var reactions []github.Reactions
	if err := json.NewDecoder(resp.Body).Decode(&reactions); err != nil {
		return 0, 0, fmt.Errorf("failed to decode reactions JSON: %w", err)
	}

	upvotes, downvotes := 0, 0
	for _, reaction := range reactions {
		switch reaction.Content {
		case "+1":
			upvotes++
		case "-1":
			downvotes++
		}
	}

	return upvotes, downvotes, nil
}
