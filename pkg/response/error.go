package response

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error,omitempty"`
	Trace   error  `json:"trace,omitempty"`
	Message string `json:"message"`
}
