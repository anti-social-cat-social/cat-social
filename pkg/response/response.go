package response

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func GenerateResponse(message string, data any) Response {
	response := Response{
		Message: message,
		Data:    data,
	}

	return response
}
