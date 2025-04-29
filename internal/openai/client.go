package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yemiwebby/code-review-agent/config"
)

const openaiURL = "https://api.openai.com/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

const promptTemplate = `
You are an expert Golang reviewer. 

Review only valid Go (.go) files that contain actual code changes.

Ignore non-Go files such as:
- .yml, .yaml, .md, .json, .txt, Dockerfile, LICENSE, README.md, etc.

Reference:
- Effective Go: https://go.dev/doc/effective_go
- Google Go Style Guide: https://google.github.io/styleguide/go/

For each Go file:
1. Write a **brief summary** of what was changed — avoid explaining every line.
2. Provide **actionable suggestions** (if any), based on idiomatic Go.
3. Add **praise** where deserved.

Keep responses clear, focused, and avoid overly detailed breakdowns.
Only comment when there are meaningful changes.

Code Patch:
%s
`

func AnalyzeCode(diff string) (string, error) {
	apiKey := config.OpenaiApiKey
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY not set")
	}

	prompt := fmt.Sprintf(promptTemplate, diff)

	requestData := RequestBody{
		Model: "gpt-3.5-turbo", // or "gpt-4"
		Messages: []Message{
			{Role: "system", Content: "You are a helpful and concise code reviewer."},
			{Role: "user", Content: prompt},
		},
	}

	payload, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", openaiURL, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Choices) == 0 {
		return "No suggestions from AI.", nil
	}

	return res.Choices[0].Message.Content, nil
}
