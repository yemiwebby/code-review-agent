package github_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/github"
)

func TestGetPRFiles(t *testing.T) {
	mockFiles := []github.FileChange{
		{Filename: "main.go", Patch: "@@ -1 +1 @@"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🧪 Test server received request:", r.URL.Path)

		if r.URL.Path != "/repos/any/repo/pulls/1/files" {
			http.Error(w, "unexpected path: "+r.URL.Path, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockFiles)
	}))
	defer ts.Close()

	// override base URL and token temporarily
	oldBaseURL := github.GitHubBaseURL
	oldToken := config.GithubToken

	github.GitHubBaseURL = ts.URL
	os.Setenv("GITHUB_TOKEN", "dummy")

	fmt.Println("👉 GitHubBaseURL is:", github.GitHubBaseURL)

	// Restoring the original functioon, ensures that the test environment is cleaned up and other tests are not affected.
	defer func() {
		github.GitHubBaseURL = oldBaseURL
		os.Setenv("GITHUB_TOKEN", oldToken)
	}()

	// Act
	files, err := github.GetPRFiles("any/repo", 1)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 || files[0].Filename != "main.go" {
		t.Errorf("unexpected files: %+v", files)
	}

	if files[0].Patch != "@@ -1 +1 @@" {
		t.Errorf("unexpected patch: %s", files[0].Patch)
	}

}
