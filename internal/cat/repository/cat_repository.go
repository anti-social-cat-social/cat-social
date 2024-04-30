package cat

import (
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	response "1-cat-social/pkg/response"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ICatRepository interface {
	FindAll(queryParam *dto.CatRequestQueryParams) ([]*entity.Cat, *response.ErrorResponse)
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

func (repo *CatRepository) FindAll(queryParam *dto.CatRequestQueryParams) ([]*entity.Cat, *response.ErrorResponse) {
	cats := []*entity.Cat{}

	query := repo.generateFilterCatQuery(queryParam)

	err := repo.db.Select(&cats, query)
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return cats, nil
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
				Err:     "Cat not found",
				Message: "error",
			}
		}
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
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
			Err:     "Internal server error",
			Message: "error",
		}
	}

	return &cat, nil
}

func (repo *CatRepository) generateFilterCatQuery(queryParam *dto.CatRequestQueryParams) string {
	query := "SELECT * FROM cats WHERE isdeleted = false"

	if queryParam.ID != "" {
		query += fmt.Sprintf(" AND WHERE id = '%s'", queryParam.ID)
	}
	if queryParam.Race != "" {
		query += fmt.Sprintf(" AND WHERE race = '%s'", queryParam.Race)
	}
	if queryParam.Sex != "" {
		query += fmt.Sprintf(" AND WHERE sex = '%s'", queryParam.Sex)
	}
	if queryParam.HasMatched != "" {
		query += fmt.Sprintf(" AND WHERE hasmatched = %s", queryParam.HasMatched)
	}
	if queryParam.AgeInMonth != "" {
		query += fmt.Sprintf(" AND WHERE ageinmonth = %s", queryParam.AgeInMonth)
	}
	if queryParam.Owned != "" {
		query += fmt.Sprintf(" AND WHERE owned = %s", queryParam.Owned)
	}
	if queryParam.Search != "" {
		query += fmt.Sprintf(" AND WHERE name LIKE '%%%s%%'", queryParam.Search)
	}
	if queryParam.Limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", queryParam.Limit)
	}
	if queryParam.Offset != 0 {
		query += fmt.Sprintf(" OFFSET %d", queryParam.Offset)
	}

	return query
}
