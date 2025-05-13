package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/webhook"
)

func init() {
	config.LoadEnv()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "AI PR Review Agent is running")
	case "/webhook":
		webhook.Handle(w, r)
	case "/github-app":
		webhook.GithubAppHandler(w, r)
	case "/process-reaction":
		webhook.ProcessReactions(w, r)
	case "/check-reactions":
		webhook.CheckReactionsHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func main() {
	lambda.Start(func(w http.ResponseWriter, r *http.Request) error {
		Handler(w, r)
		return nil
	})
}
