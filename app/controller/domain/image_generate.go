package domain

type RequestImageGenerate struct {
	Prompt string `json:"prompt"`
}

type ResponseImageGenerate struct {
	ImageURL string `json:"image_url"`
	Content string `json:"content"`
}