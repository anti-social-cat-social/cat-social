package cat

import (
	entity "1-cat-social/internal/cat/entity"
	response "1-cat-social/pkg/response"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ICatRepository interface {
	FindById(id string) (*entity.Cat, *response.ErrorResponse)
	Update(entity entity.Cat) (*entity.Cat, *response.ErrorResponse)
}

type CatRepository struct {
	db *sqlx.DB
}

func NewCatRepository(db *sqlx.DB) ICatRepository {
	return &CatRepository{
		db: db,
	}
}

func (repo *CatRepository) IsCatExist(id string) error {
	var cat entity.Cat

	err := repo.db.Get(cat, "SELECT * FROM cats WHERE id = $1", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("cat with ID %s does not exist", id)
		}
		return err
	}
	return nil
}

func (repo *CatRepository) FindById(id string) (*entity.Cat, *response.ErrorResponse) {
	cat := &entity.Cat{}

	err := repo.db.Get(cat, "SELECT * FROM cats WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.ErrorResponse{
				Code:    404,
				Error:   "Cat not found",
				Message: "error",
			}
		}
		return nil, &response.ErrorResponse{
			Code:    500,
			Error:   "Internal Server Error",
			Message: "error",
		}
	}

	return cat, nil
}

func (repo *CatRepository) Update(cat entity.Cat) (*entity.Cat, *response.ErrorResponse) {
	query := "UPDATE cats SET name = $1, race = $2, sex = $3, ageInMonth = $4, description = $5, imageUrls = $6"
	_, err := repo.db.Exec(query, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageUrls)

	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Error:   "Internal server error",
			Message: "error",
		}
	}

	return &cat, nil
}
