package response

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func GenerateResponse(code int, message string, data any) Response {
	response := Response{
		Message: message,
		Data:    data,
	}

	return response
}
