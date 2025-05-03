# 🤖 AI Code Review Agent for Pull Requests

An open source AI powered code review agent built with Go, integrated with GitHub and designed to work with any CI/CD system. It automatically reviews `.go` files in pull requests using OpenAI, posts contextual comments, and enforces that each comment is acknowledged before merging.

Let the bots do the reviewing — so you can focus on writing great code.

## ✨ Features

- 🔍 Reviews .go files in pull requests using OpenAI

- 💬 Posts summaries, suggestions, and praise directly on PRs

- ✅ Enforces reactions (👍 or 👎) and optional code updates before merge

- 🔒 Blocks unreviewed pull requests through CI integration

- ⚙️ Lightweight, extensible, and CI/CD agnostic

## 📦 Getting Started

### Clone the repository

```bash
git clone http://github.com/yemiwebby/code-review-agent.git
git push -u origin main
```

### Set environment Variables

Create a `.env` file in the root directory of the project and add the following environment variables:

```bash
OPENAI_API_KEY=your_openai_api_key
GITHUB_TOKEN=your_github_token
```

These credentials allow the agent to interact with OpenAI and GitHub APIs.

## Run the agent locally

To run the agent, use the following command:

```bash
go run cmd/main.go
```

The agent starts on port 8080.

## 🌐 Expose the Agent (Ngrok or Hosting)

To receive webhook events from GitHub, expose your local server using Ngrok:

```bash
ngrok http 8080
```

Copy the public URL — you’ll use it for the GitHub webhook and your CI system.

You can also deploy the agent to your preferred hosting environment or container platform.

## 🔁 Add a GitHub Webhook

In your GitHub repository:

- Navigate to Settings → Webhooks → Add webhook

- Use these values:

  - Payload URL: https://your-ngrok-url/webhook

  - Content type: application/json

  - Event: Pull request

## 🧪 Integrate with Your CI/CD Pipeline

To block unreviewed pull requests, connect your CI/CD pipeline to the agent's /check-reactions endpoint.

Example (with CircleCI)
Here’s how to use the agent in a CircleCI config:

```yaml
- run:
    name: Trigger AI gatekeeper check
    command: |
      curl -s -w "%{http_code}" -o result.txt "$AI_REVIEW_AGENT_URL/check-reactions?repo=$REPO_NAME"
```

Add these environment variables in your CI environment:

- `AI_REVIEW_AGENT_URL`: Your public Ngrok or hosted URL

- `REPO_NAME`: Format username/repository

You can replicate this check in GitHub Actions, GitLab CI, or any other system that supports HTTP calls.

## 🧠 Prompt Configuration

Want to change how the agent reviews code or switch to a different language?

[Edit the prompt logic here](https://github.com/yemiwebby/code-review-agent/blob/main/internal/openai/client.go)

## 💡 How It Works

- PR is opened

- Agent analyzes changed .go files

- Comments are posted directly on the PR

- CI checks that all comments have reactions or code updates

- If all good, merge is unblocked

## 🔧 Extend the Agent

Ideas for extending this project:

- Support other languages via prompt tuning

- Diff-aware inline suggestions

- Persistent comment state

- Turn it into a GitHub App

## 🤝 Contributing

This is an open source project — contributions are welcome! Whether you're fixing bugs, adding features, or tweaking the review prompt, feel free to open a PR.

If you find it useful:

- ⭐ Star the repo

- 🔁 Share it with your team

- 👤 Follow [@yemiwebby](https://www.linkedin.com/in/yemiwebby/) for updates

## 📄 License

MIT © [yemiwebby](https://www.linkedin.com/in/yemiwebby/)
