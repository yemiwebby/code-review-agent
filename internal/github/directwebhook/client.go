package directwebhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/yemiwebby/code-review-agent/config"
)

var (
	Mu         = &sync.Mutex{}
	AIComments = map[int]*AIComment{}
)

type AIComment struct {
	ID        int
	Body      string
	File      string
	Timestamp time.Time
	Line      int
	FilePath  string
	OldPatch  string
}

func PostReviewComment(owner, repo string, prNumber int, body, file string, line int, patch string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", owner, repo, prNumber)
	payload, _ := json.Marshal(map[string]string{"body": body})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "token "+config.GithubToken)
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

	var result struct {
		ID int `json:"id"`
	}
	json.Unmarshal(respBody, &result)

	Mu.Lock()
	defer Mu.Unlock()
	AIComments[result.ID] = &AIComment{
		ID:        result.ID,
		Body:      body,
		File:      file,
		Timestamp: time.Now(),
		Line:      line,
		FilePath:  file,
		OldPatch:  patch,
	}

	fmt.Printf("Posted review comment for %s: %s (ID: %d)\n", file, body, result.ID)
	return nil
}
