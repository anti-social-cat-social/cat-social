package cat

import (
	dto "1-cat-social/internal/cat/dto"
	repo "1-cat-social/internal/cat/repository"
	"1-cat-social/pkg/validation"

	"github.com/go-playground/validator/v10"
)

func ValidateIsCatExist(id string, r repo.ICatRepository) error {
	err := r.IsCatExist(id)

	return err
}

func ValidateUpdateCatForm(input dto.CatUpdateRequestBody) error {
	validate := validator.New()
	validate.RegisterValidation("valid_name", validation.ValidNameValidator)

	err := validate.Struct(input)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			return fieldError
		}
	}
	return nil
}
