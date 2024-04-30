package user

import (
	localError "1-cat-social/pkg/error"
	"1-cat-social/pkg/hasher"
	"errors"
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

	// Check if user with given email is already exists
	existedUser, _ := u.repo.FindByEmail(dto.Email)
	if existedUser != nil {
		return nil, localError.ErrForbidden("User already exists", errors.New("user already exists"))
	}

	// Map DTO to user entity
	// This used for storing data to database
	user := User{
		Name:  dto.Name,
		Email: dto.Email,
	}

	// Generate user password
	password, errHash := hasher.HashPassword(dto.Password)
	if errHash != nil {
		return nil, localError.ErrInternalServer(errHash.Error(), errHash)
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
