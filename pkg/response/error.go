package response

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   error  `json:"error"`
	Message string `json:"message"`
}
