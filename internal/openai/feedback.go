package openai

var AdjustPromptMsg = "Note: Review quality should improve. Be more concise and offer clearer, more actionable suggestions."

func AdjustPrompt(prompt string, upvotes, downvotes int) string {
	if downvotes > upvotes {
		return prompt + "\n\n" + AdjustPromptMsg
	}
	return prompt
}
