package localError

type GlobalError struct {
	Code    int
	Message string
	Error   error
}

// Return Not found error structure with customize message and error.
func ErrNotFound(message string, err error) *GlobalError {
	baseError := GlobalError{
		Code:    404,
		Message: message,
		Error:   err,
	}

	return &baseError
}
