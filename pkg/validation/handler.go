package validation

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

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

func ValidNameValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Improved regex pattern:
	// - Allows spaces, hyphens, and apostrophes within a name part
	// - Uses a character class for allowed characters (letters, spaces, hyphens, apostrophes)
	pattern := `^[[:alpha:]]+(?: [[:alpha:]]+|[-\'][:alpha:]]+)*$`

	// COmpile regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Println("Error compiling regex")
	}

	return re.MatchString(value)
}
