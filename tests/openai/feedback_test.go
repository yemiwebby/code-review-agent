package openai

import (
	"testing"

	"github.com/yemiwebby/code-review-agent/internal/openai"
)

func TestAdjustPrompt(t *testing.T) {

	tests := []struct {
		name           string
		upvotes        int
		downvotes      int
		originalPrompt string
		expected       string
	}{
		{
			name:           "when downvotes > upvotes",
			upvotes:        1,
			downvotes:      2,
			originalPrompt: "Hey",
			expected:       "Hey\n\n" + openai.AdjustPromptMsg,
		},
		{
			name:           "when downvotes <= upvotes",
			upvotes:        1,
			downvotes:      1,
			originalPrompt: "Hey",
			expected:       "Hey",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := openai.AdjustPrompt(tc.originalPrompt, tc.upvotes, tc.downvotes)

			if result != tc.expected {
				t.Errorf("upexpected prompt.\nGot: %q\nWant: %q", result, tc.expected)
			}
		})
	}
}
