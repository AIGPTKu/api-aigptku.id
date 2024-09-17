package domain

type AskContent struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type FuncCall struct {
	Name string `json:"name"`
	Arguments Arguments `json:"arguments"`
}

type Arguments struct {
	Query string `json:"query,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}

type RequestAsk struct {
	Room string `json:"room"`
	Contents []AskContent `json:"contents"`
}

type ResponseAsk struct {
	Content string `json:"content"`
	FuncCall FuncCall `json:"function_call"`
}