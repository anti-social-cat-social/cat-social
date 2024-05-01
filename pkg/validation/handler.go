package validate

import "github.com/go-playground/validator/v10"

type validationError struct {
	Field    string `json:"field"`
	Messsage string `json:"message"`
}

func GenerateStructValidationError(errorValidate validator.ValidationErrors) []validationError {
	var result []validationError

	for _, err := range errorValidate {
		e := validationError{
			Field:    err.Field(),
			Messsage: err.Error(),
		}

		result = append(result, e)
	}

	return result
}
