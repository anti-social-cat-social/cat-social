package cat

import (
	"1-cat-social/config"
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	"1-cat-social/pkg/logger"
	response "1-cat-social/pkg/response"

	"github.com/jmoiron/sqlx"
)

type IMatchRepository interface {
	MatchCat(req *dto.CatMatchRequest, issuerID string) *response.ErrorResponse
	WithTrx(trxHandle *sqlx.Tx) *matchRepository
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
