package cat

import (
	"1-cat-social/config"
	dto "1-cat-social/internal/cat/dto"
	entity "1-cat-social/internal/cat/entity"
	localError "1-cat-social/pkg/error"
	"1-cat-social/pkg/logger"
	response "1-cat-social/pkg/response"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ICatRepository interface {
	FindAll(queryParam *dto.CatRequestQueryParams, userID string) ([]*entity.Cat, *response.ErrorResponse)
	FindById(id string) (*entity.Cat, *response.ErrorResponse)
	Update(entity entity.Cat) (*entity.Cat, *response.ErrorResponse)
	WithTrx(*sqlx.Tx) *catRepository
	IsCatExist(id string) error
	Create(entity entity.Cat) (*entity.Cat, *localError.GlobalError)
	Delete(entity entity.Cat) *localError.GlobalError
}

type catRepository struct {
	db   *sqlx.DB
	tXdb *sqlx.Tx
}

func NewCatRepository(db *sqlx.DB) ICatRepository {
	return &catRepository{
		db:   db,
		tXdb: nil,
	}
}

func (repo *catRepository) getDB() config.DB {
	if repo.tXdb != nil {
		return repo.tXdb
	}
	return repo.db
}

func (repo *catRepository) WithTrx(trxHandle *sqlx.Tx) *catRepository {
	if trxHandle == nil {
		logger.Info("Transaction Database not found")
		return repo
	}
	repo.tXdb = trxHandle
	return repo
}

func (repo *catRepository) FindAll(queryParam *dto.CatRequestQueryParams, userID string) ([]*entity.Cat, *response.ErrorResponse) {
	cats := []*entity.Cat{}

	query := repo.generateFilterCatQuery(queryParam)

	var err error
	if queryParam.Owned != "" {
		err = repo.db.Select(&cats, query, userID)
	} else {
		err = repo.db.Select(&cats, query)
	}
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal Server Error",
			Message: err.Error(),
		}
	}

	return cats, nil
}

func (repo *catRepository) IsCatExist(id string) error {
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

func (repo *catRepository) FindById(id string) (*entity.Cat, *response.ErrorResponse) {
	cat := &entity.Cat{}

	err := repo.getDB().Get(cat, "SELECT * FROM cats WHERE id = $1 and isdeleted is false", id)
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

// Store new cat data
func (repo *catRepository) Create(cat entity.Cat) (*entity.Cat, *localError.GlobalError) {
	// Generate UUID
	catId := uuid.NewString()
	cat.ID = catId

	// Craeted tim generator
	createdAt := time.Now()
	cat.CreatedAt = createdAt

	query := "INSERT INTO cats (id,name, race, sex,ageinmonth,description,imageurls,hasmatched,ownerid,createdat) values (:id,:name, :race,:sex,:ageinmonth,:description,:imageurls,:hasmatched,:ownerid,:createdat)"
	_, err := repo.getDB().NamedExec(query, &cat)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return &cat, nil
}

func (repo *catRepository) Update(cat entity.Cat) (*entity.Cat, *response.ErrorResponse) {
	var err error

	query := "UPDATE cats SET name = $1, race = $2, sex = $3, ageInMonth = $4, description = $5, imageUrls = $6, hasmatched = $7 WHERE id = $8"
	_, err = repo.getDB().Exec(query, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageUrls, cat.HasMatched, cat.ID)
	if err != nil {
		return nil, &response.ErrorResponse{
			Code:    500,
			Err:     "Internal server error",
			Message: err.Error(),
		}
	}

	return &cat, nil
}

func (repo *catRepository) Delete(cat entity.Cat) *localError.GlobalError {
	_, err := repo.getDB().Exec("UPDATE cats set isdeleted = true where id = $1", cat.ID)
	if err != nil {
		return localError.ErrInternalServer(err.Error(), err)
	}
	return nil
}

func (repo *catRepository) generateFilterCatQuery(queryParam *dto.CatRequestQueryParams) string {
	query := "SELECT * FROM cats WHERE isdeleted = false"

	query = addIDFilter(query, queryParam.ID)
	query = addRaceFilter(query, queryParam.Race)
	query = addSexFilter(query, queryParam.Sex)
	query = addHasMatchedFilter(query, queryParam.HasMatched)
	query = addAgeInMonthFilter(query, queryParam.AgeInMonth)
	query = addOwnedFilter(query, queryParam.Owned)
	query = addSearchFilter(query, queryParam.Search)

	query += " ORDER BY createdat DESC"

	query = addLimitOffset(query, queryParam.Limit, queryParam.Offset)

	return query
}

func addIDFilter(query string, id string) string {
	if id != "" {
		query += fmt.Sprintf(" AND id = '%s'", id)
	}
	return query
}

func addRaceFilter(query string, race string) string {
	if race != "" {
		query += fmt.Sprintf(" AND race = '%s'", race)
	}
	return query
}

func addSexFilter(query string, sex string) string {
	if sex != "" {
		query += fmt.Sprintf(" AND sex = '%s'", sex)
	}
	return query
}

func addHasMatchedFilter(query string, hasMatched string) string {
	if hasMatched != "" {
		query += fmt.Sprintf(" AND hasmatched = %s", hasMatched)
	}
	return query
}

func addAgeInMonthFilter(query string, ageInMonth string) string {
	if ageInMonth != "" {
		if strings.Contains(ageInMonth, ">") || strings.Contains(ageInMonth, "<") {
			query += fmt.Sprintf(" AND ageinmonth %s", ageInMonth)
		} else {
			query += fmt.Sprintf(" AND ageinmonth %s", ageInMonth)
		}
	}
	return query
}

func addOwnedFilter(query string, owned string) string {
	if owned != "" {
		if owned == "true" {
			query += " AND ownerid = $1"
		} else {
			query += " AND ownerid != $1"
		}
	}
	return query
}

func addSearchFilter(query string, search string) string {
	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", search)
	}
	return query
}

func addLimitOffset(query string, limit, offset int) string {
	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	} else {
		query += " LIMIT 10"
	}
	if offset != 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}
	return query
}
