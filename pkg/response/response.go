package response

type Response struct {
	Message interface{} `json:"message"`
	Data    any         `json:"data,omitempty"`
}

func GenerateResponse(message interface{}, data any) Response {
	response := Response{
		Message: message,
		Data:    data,
	}

	return response
}
