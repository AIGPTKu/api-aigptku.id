package domain

type AskContent struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type RequestAsk struct {
	Room string `json:"room"`
	Contents []AskContent `json:"contents"`
}

type ResponseAsk struct {
	Content string `json:"content"`
}