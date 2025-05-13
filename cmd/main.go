package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/webhook"
)

func main() {
	config.LoadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "AI PR Review Agent is running")
	})
	mux.HandleFunc("/webhook", webhook.Handle)
	mux.HandleFunc("/github-app", webhook.GithubAppHandler)
	mux.HandleFunc("/process-reaction", webhook.ProcessReactions)
	mux.HandleFunc("/check-reactions", webhook.CheckReactionsHandler)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
