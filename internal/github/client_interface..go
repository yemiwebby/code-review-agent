package github

type GithubClientInterface interface {
	PostReviewComment(owner, repo string, prNumber int, body, file, commitID string, line int, patch string) (int, error)
	FetchReactions(owner, repo string, commentID int) (int, int, error)
}
