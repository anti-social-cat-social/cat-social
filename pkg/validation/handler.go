package validation

import "github.com/go-playground/validator/v10"

type validationError struct {
	Field    string `json:"field"`
	Messsage string `json:"message"`
}

func GenerateStructValidationError(err error) []validationError {
	var result []validationError

	for _, err := range err.(validator.ValidationErrors) {
		e := validationError{
			Field:    err.Field(),
			Messsage: err.Error(),
		}

		result = append(result, e)
	}

	return result
}
