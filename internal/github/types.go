package github

type FileChange struct {
	SHA      string `json:"sha"`
	Filename string `json:"filename"`
	Patch    string `json:"patch"`
}
