package cat

import (
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	repo "1-cat-social/internal/cat/repository"
	localError "1-cat-social/pkg/error"
	"1-cat-social/pkg/response"
	"errors"
)

type ICatUsecase interface {
	GetAll(queryParam *dto.CatRequestQueryParams, userID string) ([]*entity.Cat, *response.ErrorResponse)
	Update(id string, dto dto.CatUpdateRequestBody) (*entity.Cat, *response.ErrorResponse)
	Create(dto dto.CatUpdateRequestBody, userID string) (*entity.Cat, *localError.GlobalError)
	Delete(id string, userId string) *localError.GlobalError
}

type CatUsecase struct {
	repo repo.ICatRepository
}

func NewCatUsecase(repo repo.ICatRepository) ICatUsecase {
	return &CatUsecase{
		repo: repo,
	}
}

func (uc *CatUsecase) GetAll(queryParam *dto.CatRequestQueryParams, userID string) ([]*entity.Cat, *response.ErrorResponse) {
	return uc.repo.FindAll(queryParam, userID)
}

func (uc *CatUsecase) Create(dto dto.CatUpdateRequestBody, userID string) (*entity.Cat, *localError.GlobalError) {
	catData := entity.Cat{
		Name:        dto.Name,
		Race:        dto.Race,
		Sex:         dto.Sex,
		AgeInMonth:  dto.AgeInMonth,
		Description: dto.Description,
		OwnerId:     userID,
		ImageUrls:   dto.ImageUrls,
	}

	cat, err := uc.repo.Create(catData)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (uc *CatUsecase) Update(id string, input dto.CatUpdateRequestBody) (*entity.Cat, *response.ErrorResponse) {
	cat, error := uc.repo.FindById(id)
	if error != nil {
		return nil, error
	}

	if cat.Sex != input.Sex && cat.HasMatched {
		return nil, &response.ErrorResponse{
			Code:    400,
			Err:     "Can't update cat that already matched",
			Message: "error",
		}
	}

	cat.Name = input.Name
	cat.Race = input.Race
	cat.Sex = input.Sex
	cat.AgeInMonth = input.AgeInMonth
	cat.Description = input.Description
	cat.ImageUrls = input.ImageUrls

	return uc.repo.Update(*cat)
}

func (uc *CatUsecase) Delete(id string, userID string) *localError.GlobalError {
	cat, errpr := uc.repo.FindById(id)
	if errpr != nil {
		return &localError.GlobalError{
			Code:    404,
			Message: errpr.Err,
			Error:   errpr.Trace,
		}
	}

	// Check if the cat is on the right OwnerId
	if cat.OwnerId != userID {
		return localError.ErrForbidden("Owner cat invalid", errors.New("owner cat invalid"))
	}

	// Delete cat
	if err := uc.repo.Delete(*cat); err != nil {
		return err
	}

	return nil
}
