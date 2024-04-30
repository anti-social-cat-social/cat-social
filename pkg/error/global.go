package localError

import "net/http"

type GlobalError struct {
	Code    int
	Message string
	Error   error
}

// Return Not found error structure with customize message and error.
func ErrBase(code int, message string, err error) *GlobalError {
	baseError := GlobalError{
		Code:    code,
		Message: message,
		Error:   err,
	}

	return &baseError
}

// Return internal server error structure with customize message and error.
func ErrInternalServer(message string, err error) *GlobalError {
	baseError := ErrBase(http.StatusInternalServerError, message, err)

	return baseError
}

// Return unauthorized structure with customize message and error.
func ErrUnauthorized(message string, err error) *GlobalError {
	baseError := ErrBase(http.StatusUnauthorized, message, err)

	return baseError
}

// Return unauthorized structure with customize message and error.
func ErrForbidden(message string, err error) *GlobalError {
	baseError := ErrBase(http.StatusForbidden, message, err)

	return baseError
}

// Return Not found error structure with customize message and error.
func ErrNotFound(message string, err error) *GlobalError {
	baseError := ErrBase(http.StatusNotFound, message, err)

	return baseError
}
