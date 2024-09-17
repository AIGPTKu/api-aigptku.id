package model

type ResponseDallE3 struct {
	Created int64 `json:"created"`
	Data []struct{
		RevisedPrompt string `json:"revised_prompt"`
		Url string `json:"url"`
	} `json:"data"`
	Error struct {
		Code string `json:"code"`
		Message string `json:"message"`
		Type string `json:"type"`
	} `json:"error"`
}