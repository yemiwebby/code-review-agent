package github

type GithubClientInterface interface {
	PostReviewComment(owner, repo string, prNumber int, body, file, commitID string, line int, patch string) error
}
