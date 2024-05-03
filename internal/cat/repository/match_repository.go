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
	MatchCat(req *dto.CatMatchRequest, issuerID, targetOwnerID string) *response.ErrorResponse
	FindById(id string) (*entity.Match, *response.ErrorResponse)
	IsCatAlreadyMatched(userCatId, targetCatId string) (bool, *response.ErrorResponse)
	FindByCatId(id string) (entity.Match, *response.ErrorResponse)
	DeleteMatch(issuerCatID, targetCatID string, matchID string) *response.ErrorResponse
	ApproveMatch(matchID string) *response.ErrorResponse
	RejectMatch(matchID string) *response.ErrorResponse
	WithTrx(trxHandle *sqlx.Tx) *matchRepository
	RemoveTrx()
	FindAllByUserId(userId string) ([]dto.MatchResponse, *response.ErrorResponse)
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

func (repo *matchRepository) MatchCat(req *dto.CatMatchRequest, issuerID, targetOwnerId string) *response.ErrorResponse {
	payload := map[string]interface{}{
		"usercatid":        req.UserCatId,
		"matchcatid":       req.MatchCatId,
		"message":          req.Message,
		"status":           entity.Submitted,
		"issuedby":         issuerID,
		"target_cat_owner": targetOwnerId,
	}

	res, err := repo.getDB().NamedExec(`INSERT INTO matches (issuer_cat_id, target_cat_id, message, status, issuedby, target_cat_owner) VALUES (:usercatid,:matchcatid,:message,:status,:issuedby,:target_cat_owner) RETURNING id`, payload)
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

func (repo *matchRepository) FindAllByUserId(userId string) ([]dto.MatchResponse, *response.ErrorResponse) {
	matches := []dto.MatchResponse{}

	query := `
        SELECT
		m.id AS id,
		iss.name AS "issuedBy.name",
		iss.email AS "issuedBy.email",
		iss.createdat AS "issuedBy.createdat",
		catTar.id AS "matchCatDetail.id",
		catTar.name AS "matchCatDetail.name",
		catTar.race AS "matchCatDetail.race",
		catTar.sex AS "matchCatDetail.sex",
		catTar.description AS "matchCatDetail.description",
		catTar.ageinmonth AS "matchCatDetail.ageinmonth",
		catTar.imageurls AS "matchCatDetail.imageurls",
		catTar.hasmatched AS "matchCatDetail.hasmatched",
		catTar.createdat AS "matchCatDetail.createdat",
		catIss.id AS "userCatDetail.id",
		catIss.name AS "userCatDetail.name",
		catIss.race AS "userCatDetail.race",
		catIss.sex AS "userCatDetail.sex",
		catIss.description AS "userCatDetail.description",
		catIss.ageinmonth AS "userCatDetail.ageinmonth",
		catIss.imageurls AS "userCatDetail.imageurls",
		catIss.hasmatched AS "userCatDetail.hasmatched",
		catIss.createdat AS "userCatDetail.createdat",
		m.message AS message,
		m.createdat AS createdat
	FROM
		matches m
		JOIN cats catIss ON m.issuer_cat_id = catIss.id
		JOIN users iss ON m.issuedby = iss.id
		JOIN cats catTar ON m.target_cat_id = catTar.id
		JOIN users tar ON m.target_cat_owner = tar.id
	WHERE
		m.issuedby = $1
		or m.target_cat_owner = $1;
    `
	// Queryx is used here as it handles scanning into structs automatically
	err := repo.db.Select(&matches, query, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return matches, nil
}
