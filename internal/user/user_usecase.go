package user

import (
	localError "1-cat-social/pkg/error"
	"1-cat-social/pkg/hasher"

	"github.com/google/uuid"
)

type IUserUsecase interface {
	FindByEmail(email string) (*User, *localError.GlobalError)
	Create(dto UserDTO) (*User, *localError.GlobalError)
}

type userUsecase struct {
	repo IUserRepository
}

// Create implements IUserUsecase.
func (u *userUsecase) Create(dto UserDTO) (*User, *localError.GlobalError) {
	// Validate user request first

	// Map DTO to user entity
	// This used for storing data to database
	user := User{
		Name:  dto.Name,
		Email: dto.Email,
	}

	// Generate user UUID
	userId := uuid.NewString()
	user.ID = userId

	// Generate user password
	password, err := hasher.HashPassword(dto.Password)
	if err != nil {
		return nil, nil
	}
	// Assign user password to struct if not error
	user.Password = password

	return u.repo.Create(user)
}

// FindByEmail implements IUserUsecase.
func (u *userUsecase) FindByEmail(email string) (*User, *localError.GlobalError) {
	return u.repo.FindByEmail(email)
}

func NewUserUsecase(repo IUserRepository) IUserUsecase {
	return &userUsecase{
		repo: repo,
	}
}
