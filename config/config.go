package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	GithubToken  string
	OpenaiApiKey string
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".env not loaded:", err)
	}

	GithubToken = os.Getenv("GITHUB_TOKEN")
	OpenaiApiKey = os.Getenv("OPENAI_API_KEY")

}
