package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (a *AppAuthenticator) GetInstallationToken(installationID int64) (string, error) {
	jwtToken, err := a.GenerateJWT()
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID)
	fmt.Printf("Requesting installation token from URL: %s\n", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get installation token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Printf("GitHub API Response: %s\n", string(body))

	// if resp.StatusCode != http.StatusCreated {
	// 	return "", fmt.Errorf("failed to get installation token, status code: %d", resp.StatusCode)
	// }

	var result struct {
		Token string `json:"token"`
	}

	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return "", fmt.Errorf("failed to parse response body: %w", err)
	// }

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response body: %v", err)
	}

	return result.Token, nil
}
