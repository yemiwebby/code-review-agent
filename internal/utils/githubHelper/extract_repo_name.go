package githubhelper

import (
	"fmt"
	"strings"
)

func ExtractRepoName(owner, repo string) (string, error) {
	// extract the actual repo name (without the owner prefix)
	slashIndex := len(owner) + 1
	if slashIndex < len(repo) && repo[slashIndex-1] == '/' {
		return repo[slashIndex:], nil
	}

	return "", fmt.Errorf("invalid repository format: %s/%s", owner, repo)
}

// SplitRepoName extracts the owner and repository name from the full name
func SplitRepoName(repo string) (string, string, error) {
	parts := strings.Split(repo, "/")
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("invalid repository format: %s", repo)
}
