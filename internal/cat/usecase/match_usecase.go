package cat

import (
	dto "1-cat-social/internal/cat/dto"
	repo "1-cat-social/internal/cat/repository"
	"1-cat-social/pkg/response"
)

type IMatchUsecase interface {
	Match(req *dto.CatMatchRequest, userID string) *response.ErrorResponse
}

type MatchUsecase struct {
	catRepository   repo.ICatRepository
	matchRepository repo.IMatchRepository
}

func NewMatchUsecase(cr repo.ICatRepository, mr repo.IMatchRepository) IMatchUsecase {
	return &MatchUsecase{
		catRepository:   cr,
		matchRepository: mr,
	}
}

func (uc *MatchUsecase) Match(req *dto.CatMatchRequest, userID string) *response.ErrorResponse {
	// check if matchCatId and userCatId is exist
	matchCat, err := uc.catRepository.FindById(req.MatchCatId)
	if err != nil {
		return &response.ErrorResponse{
			Code:    404,
			Err:     "Match cat not found",
			Message: "error",
		}
	}
	userCat, err := uc.catRepository.FindById(req.UserCatId)
	if err != nil {
		return &response.ErrorResponse{
			Code:    404,
			Err:     "User cat not found",
			Message: "error",
		}
	}

	// check if userCatId is belongs to user
	if userCat.OwnerId != userID {
		return &response.ErrorResponse{
			Code:    404,
			Err:     "Cat is not belongs to you",
			Message: "error",
		}
	}
	// check if matchCatId and userCatId is not from the same owner
	if matchCat.OwnerId == userID {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Can't match with your own cat",
			Message: "error",
		}
	}

	// check if cats gender is the same or not
	if matchCat.Sex == userCat.Sex {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Can't match with the same gender",
			Message: "error",
		}
	}

	// check if cats already matched or not
	if matchCat.HasMatched || userCat.HasMatched {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Cat already matched",
			Message: "error",
		}
	}

	// update hasMatched to true
	matchCat.HasMatched = true
	userCat.HasMatched = true

	matched, err := uc.matchRepository.MatchCat(userCat, matchCat, req.Message)
	if err != nil {
		return err
	}

	if matched.ID == "" {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: "failed to match cat",
		}
	}

	return nil
}
