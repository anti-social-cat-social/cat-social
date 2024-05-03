package cat

import (
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	repo "1-cat-social/internal/cat/repository"
	"1-cat-social/pkg/response"

	"github.com/jmoiron/sqlx"
)

type IMatchUsecase interface {
	WithTrx(*sqlx.Tx) *matchUsecase
	Match(req *dto.CatMatchRequest, userID string) *response.ErrorResponse
	Approve(req *dto.MatchApproveRequest, userID string) *response.ErrorResponse
	Reject(req *dto.MatchApproveRequest, userID string) *response.ErrorResponse
	GetMatches(userID string) ([]dto.MatchResponse, *response.ErrorResponse)
}

type matchUsecase struct {
	catRepository   repo.ICatRepository
	matchRepository repo.IMatchRepository
}

func NewMatchUsecase(cr repo.ICatRepository, mr repo.IMatchRepository) IMatchUsecase {
	return &matchUsecase{
		catRepository:   cr,
		matchRepository: mr,
	}
}

func (uc *matchUsecase) WithTrx(trxHandle *sqlx.Tx) *matchUsecase {
	uc.catRepository = uc.catRepository.WithTrx(trxHandle)
	uc.matchRepository = uc.matchRepository.WithTrx(trxHandle)
	return uc
}

func (uc *matchUsecase) Match(req *dto.CatMatchRequest, userID string) *response.ErrorResponse {
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

	isMatched, err := uc.matchRepository.IsCatAlreadyMatched(req.UserCatId, req.MatchCatId)
	if err != nil {
		return err
	}

	if isMatched {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Cat already requested to match",
			Message: "error",
		}
	}

	err = uc.matchRepository.MatchCat(req, userCat.OwnerId, matchCat.OwnerId)
	if err != nil {
		return err
	}

	uc.catRepository.RemoveTrx()
	uc.matchRepository.RemoveTrx()

	return nil
}

func (uc *matchUsecase) Approve(req *dto.MatchApproveRequest, userID string) *response.ErrorResponse {
	// check if matchId is exist
	match, err := uc.matchRepository.FindById(req.MatchId)
	if err != nil {
		return err
	}

	if match.IsDeleted {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Match is no longer valid",
			Message: "error",
		}
	}

	if match.Status != entity.Submitted {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Match is no longer valid",
			Message: "error",
		}
	}

	targetCat, err := uc.catRepository.FindById(match.TargetCatId)
	if err != nil {
		return err
	}

	if targetCat.OwnerId != userID {
		return &response.ErrorResponse{
			Code:    403,
			Err:     "You are not the owner of the cat",
			Message: "error",
		}
	}

	err = uc.matchRepository.ApproveMatch(req.MatchId)
	if err != nil {
		return err
	}

	err = uc.matchRepository.DeleteMatch(match.IssuerCatId, match.TargetCatId, req.MatchId)
	if err != nil {
		return err
	}

	issuerCat, err := uc.catRepository.FindById(match.IssuerCatId)
	if err != nil {
		return err
	}
	// update hasMatched to true
	issuerCat.HasMatched = true
	targetCat.HasMatched = true

	_, err = uc.catRepository.Update(*issuerCat)
	if err != nil {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	_, err = uc.catRepository.Update(*targetCat)
	if err != nil {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	uc.catRepository.RemoveTrx()
	uc.matchRepository.RemoveTrx()

	return nil
}

func (uc *matchUsecase) Reject(req *dto.MatchApproveRequest, userID string) *response.ErrorResponse {
	// check if matchId is exist
	match, err := uc.matchRepository.FindById(req.MatchId)
	if err != nil {
		return err
	}

	if match.IsDeleted {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Match is no longer valid",
			Message: "error",
		}
	}

	if match.Status != entity.Submitted {
		return &response.ErrorResponse{
			Code:    400,
			Err:     "Match is no longer valid",
			Message: "error",
		}
	}

	targetCat, err := uc.catRepository.FindById(match.TargetCatId)
	if err != nil {
		return err
	}

	if targetCat.OwnerId != userID {
		return &response.ErrorResponse{
			Code:    403,
			Err:     "You are not the owner of the cat",
			Message: "error",
		}
	}

	err = uc.matchRepository.RejectMatch(req.MatchId)
	if err != nil {
		return err
	}

	uc.catRepository.RemoveTrx()
	uc.matchRepository.RemoveTrx()

	return nil
}

func (uc *matchUsecase) GetMatches(userID string) ([]dto.MatchResponse, *response.ErrorResponse) {
	matches, err := uc.matchRepository.FindAllByUserId(userID)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
