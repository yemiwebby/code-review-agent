# AI Pull Request Review Agent

## Getting Started

This project is designed to assist developers in reviewing pull requests using AI. The agent can analyze code changes, provide feedback, and suggest improvements. It is built using Golang and leverages the OpenAI API for natural language processing.

## Clone the repository

```bash
git clone http://github.com/yemiwebby/code-review-agent.git
git push -u origin main
```

## Environment Variables

Create a `.env` file in the root directory of the project and add the following environment variables:

```bash
OPENAI_API_KEY=your_openai_api_key
GITHUB_TOKEN=your_github_token
```

## Run the Agent

To run the agent, use the following command:

```bash
go run cmd/main.go
```
