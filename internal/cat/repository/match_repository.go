package cat

import (
	"1-cat-social/config"
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	"1-cat-social/pkg/logger"
	response "1-cat-social/pkg/response"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type IMatchRepository interface {
	MatchCat(req *dto.CatMatchRequest, issuerID string) *response.ErrorResponse
	FindById(id string) (*entity.Match, *response.ErrorResponse)
	IsCatAlreadyMatched(userCatId, targetCatId string) (bool, *response.ErrorResponse)
	FindByCatId(id string) (entity.Match, *response.ErrorResponse)
	DeleteMatch(issuerCatID, targetCatID string, matchID string) *response.ErrorResponse
	ApproveMatch(matchID string) *response.ErrorResponse
	RejectMatch(matchID string) *response.ErrorResponse
	WithTrx(trxHandle *sqlx.Tx) *matchRepository
	RemoveTrx()
}

type matchRepository struct {
	db   *sqlx.DB
	tXdb *sqlx.Tx
}

func NewMatchRepository(db *sqlx.DB) IMatchRepository {
	return &matchRepository{
		db:   db,
		tXdb: nil,
	}
}

func (repo *matchRepository) getDB() config.DB {
	if repo.tXdb != nil {
		return repo.tXdb
	}
	return repo.db
}

func (repo *matchRepository) WithTrx(trxHandle *sqlx.Tx) *matchRepository {
	if trxHandle == nil {
		logger.Info("Transaction Database not found")
		return repo
	}
	repo.tXdb = trxHandle
	return repo
}

func (repo *matchRepository) RemoveTrx() {
	repo.tXdb = nil
}

func (repo *matchRepository) MatchCat(req *dto.CatMatchRequest, issuerID string) *response.ErrorResponse {
	payload := map[string]interface{}{
		"usercatid":  req.UserCatId,
		"matchcatid": req.MatchCatId,
		"message":    req.Message,
		"status":     entity.Submitted,
		"issuedby":   issuerID,
	}

	res, err := repo.getDB().NamedExec(`INSERT INTO matches (issuer_cat_id, target_cat_id, message, status, issuedby) VALUES (:usercatid,:matchcatid,:message,:status,:issuedby) RETURNING id`, payload)
	if err != nil {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	if numRows, err := res.RowsAffected(); err != nil && numRows == 0 {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: "Failed to insert match",
		}
	}

	return nil
}

func (repo *matchRepository) FindById(id string) (*entity.Match, *response.ErrorResponse) {
	match := entity.Match{}

	err := repo.getDB().Get(&match, "SELECT * FROM matches WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Code:    404,
				Err:     "Match not found",
				Message: err.Error(),
			}
		}
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return &match, nil
}

func (repo *matchRepository) IsCatAlreadyMatched(userCatId, targetCatId string) (bool, *response.ErrorResponse) {
	match := entity.Match{}

	err := repo.getDB().Get(&match, `SELECT * FROM matches WHERE (issuer_cat_id = $1 AND target_cat_id = $2)`, userCatId, targetCatId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, &response.ErrorResponse{
				Code:    500,
				Err:     "Internal Server Error",
				Message: err.Error(),
			}
		}
	}

	if match.ID == "" {
		return false, nil
	}

	return true, nil
}

func (repo *matchRepository) FindByCatId(id string) (entity.Match, *response.ErrorResponse) {
	match := entity.Match{}

	err := repo.getDB().Get(&match, `SELECT * FROM matches WHERE (issuer_cat_id = $1 OR target_cat_id = $1)`, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return match, &response.ErrorResponse{
				Code:    500,
				Err:     "Internal Server Error",
				Message: err.Error(),
			}
		}
	}

	return match, nil
}

func (repo *matchRepository) DeleteMatch(issuerCatID, targetCatID string, matchID string) *response.ErrorResponse {
	_, err := repo.getDB().Exec(`UPDATE matches SET isdeleted = true WHERE id != $3 AND (issuer_cat_id = $1 OR target_cat_id = $1 OR issuer_cat_id = $2 OR target_cat_id = $2)`, issuerCatID, targetCatID, matchID)
	if err != nil {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return nil
}

func (repo *matchRepository) ApproveMatch(matchID string) *response.ErrorResponse {
	_, err := repo.getDB().Exec(`UPDATE matches SET status = $1 WHERE id = $2`, entity.Approved, matchID)
	if err != nil {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return nil
}

func (repo *matchRepository) RejectMatch(matchID string) *response.ErrorResponse {
	_, err := repo.getDB().Exec(`UPDATE matches SET status = $1 WHERE id = $2`, entity.Rejected, matchID)
	if err != nil {
		return &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return nil
}
