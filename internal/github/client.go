package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yemiwebby/code-review-agent/config"
)

type FileChange struct {
	Filename string `json:"filename"`
	Patch    string `json:"patch"`
}

var GitHubBaseURL = "https://api.github.com"

func GetPRFiles(repo string, prNumber int) ([]FileChange, error) {
	url := fmt.Sprintf("%s/repos/%s/pulls/%d/files", GitHubBaseURL, repo, prNumber)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "token "+config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("GitHub response body:", string(bodyBytes))

	var files []FileChange
	if err := json.Unmarshal(bodyBytes, &files); err != nil {
		return nil, err
	}

	return files, nil
}
