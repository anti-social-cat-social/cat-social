package auth

import (
	"1-cat-social/internal/user"
	localError "1-cat-social/pkg/error"
	"1-cat-social/pkg/hasher"
	localJwt "1-cat-social/pkg/jwt"
	"time"
)

type IAuthUsecase interface {
	Login(dto user.LoginDTO) (*authResponse, *localError.GlobalError)
	Register(dto user.UserDTO) (*authResponse, *localError.GlobalError)
}

type authUsecase struct {
	userUc user.IUserUsecase
}

type authResponse struct {
	UserID    string `json:"userId"`
	Token     string `json:"accessToken"`
	ExpiredAt string `json:"expiredAt"`
}

var tokenExpirationHour time.Duration = time.Duration(8)

// Login implements IAuthUsecase.
func (a *authUsecase) Login(dto user.LoginDTO) (*authResponse, *localError.GlobalError) {
	// Search user by email
	result, errUser := a.userUc.FindByEmail(dto.Email)
	if errUser != nil {
		return nil, errUser
	}

	// Compare user password with stored password
	err := hasher.CheckPassword(result.Password, dto.Password)
	if err != nil {
		return nil, localError.ErrUnauthorized("Credential not valid", err)
	}

	// Generate token if no error happened above
	// Token generated using JWT scheme
	token, err := localJwt.GenerateToken(*result)
	if err != nil {
		return nil, localError.ErrInternalServer(err.Error(), err)
	}

	// Map claim to auth response auth
	response := authResponse{
		UserID:    result.ID,
		Token:     token,
		ExpiredAt: time.Now().Add(8 * time.Hour).Format(time.RFC3339),
	}

	return &response, nil
}

// Register implements IAuthUsecase.
func (a *authUsecase) Register(dto user.UserDTO) (*authResponse, *localError.GlobalError) {
	// Create user data and generate token
	user, err := a.userUc.Create(dto)
	if err != nil {
		return nil, err
	}

	// Generate token if user successfully created
	token, errToken := localJwt.GenerateToken(*user)
	if errToken != nil {
		return nil, localError.ErrInternalServer(errToken.Error(), errToken)
	}

	// Map token to auth response auth
	response := authResponse{
		UserID:    user.Name,
		Token:     token,
		ExpiredAt: time.Now().Add(8 * time.Hour).Format(time.RFC3339),
	}

	return &response, nil
}

func NewAuthUsecase(userUc user.IUserUsecase) IAuthUsecase {
	return &authUsecase{
		userUc: userUc,
	}
}
