package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yemiwebby/code-review-agent/config"
)

type Reactions struct {
	Content string `json:"content"`
}

func FetchReactions(repo string, commentID int) (int, int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/issues/comments/%d/reactions", repo, commentID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.squirrel-girl-preview+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var reactions []Reactions
	if err := json.Unmarshal(body, &reactions); err != nil {
		fmt.Println("⚠️ Raw GitHub response (unexpected structure):", string(body)) // Debugging aid
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
