package response

type BaseResponse struct {
	StatusCode int `json:"status_code"`
	ErrorMessage string `json:"error_message,omitempty"`
	Data any `json:"data"`
}

type ErrorPayload struct {
	
}