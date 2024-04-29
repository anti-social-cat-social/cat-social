package auth

import (
	"1-cat-social/internal/user"
	localError "1-cat-social/pkg/error"
)

type IAuthUsecase interface {
	Login() (*user.User, *localError.GlobalError)
	Register()
}

type authUsecase struct {
	userUc user.IUserUsecase
}

// Login implements IAuthUsecase.
func (a *authUsecase) Login() (*user.User, *localError.GlobalError) {
	email := "bomsiwor@gmail.com"

	result, err := a.userUc.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Register implements IAuthUsecase.
func (a *authUsecase) Register() {
	panic("unimplemented")
}

func NewAuthUsecase(userUc user.IUserUsecase) IAuthUsecase {
	return &authUsecase{
		userUc: userUc,
	}
}
