package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yemiwebby/code-review-agent/config"
	"github.com/yemiwebby/code-review-agent/internal/webhook"
)

func main() {
	config.LoadEnv()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "AI PR Review Agent is running")
	})
	mux.HandleFunc("/webhook", webhook.Handle)
	mux.HandleFunc("/process-reaction", webhook.ProcessReactions)
	mux.HandleFunc("/check-reactions", webhook.CheckReactionsHandler)

	log.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
