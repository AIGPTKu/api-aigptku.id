package model

// Define the structure of the JSON response
type GPTChoice struct {
    Delta struct {
        Content string `json:"content"`
    } `json:"delta"`
}

type ResponseGPT struct {
    Choices []GPTChoice `json:"choices"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens int `json:"total_tokens"`
    } `json:"usage"`
}