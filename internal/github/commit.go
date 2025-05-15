package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CommitInfo struct {
	SHA string `json:"sha"`
}

func GetCommitID(owner, repo string, prNumber int, token string, isApp bool) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if isApp {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.Header.Set("Authorization", "token "+token)
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get commit ID, status code: %d", resp.StatusCode)
	}

	var commitInfo struct {
		Head CommitInfo `json:"head"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&commitInfo); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return commitInfo.Head.SHA, nil
}
