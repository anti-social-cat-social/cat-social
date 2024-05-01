package cat

import (
	entity "1-cat-social/internal/cat/entity"
	"1-cat-social/pkg/logger"
	response "1-cat-social/pkg/response"
	"context"

	"github.com/jmoiron/sqlx"
)

type IMatchRepository interface {
	MatchCat(userCat *entity.Cat, matchCat *entity.Cat, msg string) (*entity.Match, *response.ErrorResponse)
}

type MatchRepository struct {
	db *sqlx.DB
}

func NewMatchRepository(db *sqlx.DB) IMatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (repo *MatchRepository) MatchCat(userCat *entity.Cat, matchCat *entity.Cat, msg string) (*entity.Match, *response.ErrorResponse) {
	match := &entity.Match{}

	// update usercar and matchcat
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}
	defer func() {
		if p := recover(); p != nil {
			// Rollback the transaction if panic occurs
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// Rollback the transaction if there's an error
			tx.Rollback()
		} else {
			// Commit the transaction if successful
			err = tx.Commit()
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	_, err = tx.Exec("UPDATE cats SET hasmatched = true WHERE id = $1", userCat.ID)
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	_, err = tx.Exec("UPDATE cats SET hasmatched = true WHERE id = $1", matchCat.ID)
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	// insert match
	err = tx.QueryRowContext(context.Background(), "INSERT INTO matches (issuer_cat_id, target_cat_id, message, status, issuedby) VALUES ($1, $2, $3, $4, $5) RETURNING id", userCat.ID, matchCat.ID, msg, entity.Submitted, userCat.OwnerId).Scan(&match.ID)
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return match, nil
}
