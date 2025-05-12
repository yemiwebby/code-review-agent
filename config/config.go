package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	GithubToken, OpenaiApiKey, GithubAppId, GithubPrivateKey, GithubWebhookSecret, OwnerUsername, Repo, PrNumber string
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env not loaded:", err)
	}

	GithubToken = os.Getenv("GITHUB_TOKEN")
	OpenaiApiKey = os.Getenv("OPENAI_API_KEY")
	GithubAppId = os.Getenv("GITHUB_APP_ID")
	GithubPrivateKey = os.Getenv("GITHUB_PRIVATE_KEY")
	GithubWebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	OwnerUsername = os.Getenv("OWNER_USERNAME")
	Repo = os.Getenv("REPO")
	PrNumber = os.Getenv("PR_NUMBER")
}
