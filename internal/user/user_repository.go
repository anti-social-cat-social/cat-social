package user

import (
	localError "1-cat-social/pkg/error"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	FindByEmail(email string) (*User, *localError.GlobalError)
	Create(entity User) (*User, *localError.GlobalError)
}

type userRepository struct {
	db *sqlx.DB
}

// Store new ueser to database
func (u *userRepository) Create(entity User) (*User, *localError.GlobalError) {
	// Generate user UUID
	userId := uuid.NewString()
	entity.ID = userId

	// Insert into database
	query := "INSERT INTO users (id, name, email, password) values (:id, :name, :email, :password)"
	_, err := u.db.NamedExec(query, &entity)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	return &entity, nil
}

// Find user by email
// This can be use for authentication process
func (u *userRepository) FindByEmail(email string) (*User, *localError.GlobalError) {
	var user User

	if err := u.db.Get(&user, "SELECT * FROM users where email=$1", email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, localError.ErrNotFound("User data not found", err)
		}

		return nil, &localError.GlobalError{
			Code:    400,
			Message: "Not found",
			Error:   err,
		}

	}

	return &user, nil
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}
