package cat

import (
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	repo "1-cat-social/internal/cat/repository"
	"1-cat-social/pkg/response"
)

type ICatUsecase interface {
	GetAll(queryParam *dto.CatRequestQueryParams) ([]*entity.Cat, *response.ErrorResponse)
	Update(id string, dto dto.CatUpdateRequestBody) (*entity.Cat, *response.ErrorResponse)
}

type CatUsecase struct {
	repo repo.ICatRepository
}

func NewCatUsecase(repo repo.ICatRepository) ICatUsecase {
	return &CatUsecase{
		repo: repo,
	}
}

func (uc *CatUsecase) GetAll(queryParam *dto.CatRequestQueryParams) ([]*entity.Cat, *response.ErrorResponse) {
	return uc.repo.FindAll(queryParam)
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
