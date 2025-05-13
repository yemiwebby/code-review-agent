package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	GithubToken         string
	OpenaiApiKey        string
	GithubAppId         string
	GithubPrivateKey    string
	GithubWebhookSecret string
	OwnerUsername       string
	Repo                string
	PrNumber            string
)

func LoadEnv() {
	// Only try to load .env file if it exists (local development)
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println(".env not loaded:", err)
		} else {
			log.Println(".env loaded successfully")
		}
	} else {
		log.Println(".env not found, skipping")
	}

	// Load environment variables (Heroku and local)
	GithubToken = os.Getenv("GITHUB_TOKEN")
	OpenaiApiKey = os.Getenv("OPENAI_API_KEY")
	GithubAppId = os.Getenv("GITHUB_APP_ID")
	GithubPrivateKey = os.Getenv("GITHUB_PRIVATE_KEY")
	GithubWebhookSecret = os.Getenv("GITHUB_WEBHOOK_SECRET")
	OwnerUsername = os.Getenv("OWNER_USERNAME")
	Repo = os.Getenv("REPO")
	PrNumber = os.Getenv("PR_NUMBER")

	// Log for visibility
	log.Printf("Environment variables loaded: GITHUB_APP_ID=%s, OWNER_USERNAME=%s, REPO=%s, PR_NUMBER=%s",
		GithubAppId, OwnerUsername, Repo, PrNumber)
}
