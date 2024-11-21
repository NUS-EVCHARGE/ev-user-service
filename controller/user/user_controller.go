package user

import (
	"github.com/NUS-EVCHARGE/ev-user-service/controller/authentication"
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
)

type UserController interface {
	GetUserInfo(jwtToken string) (dto.User, error)
}

type UserControllerImpl struct {
}

var (
	UserControllerObj UserController
)

func (u *UserControllerImpl) GetUserInfo(token string) (dto.User, error) {
	// Cast the claims
	user, err := authentication.AuthenticationControllerObj.GetUserInfo(token)

	if err != nil {
		return dto.User{}, err
	}

	var username, email string
	for _, attr := range user.UserAttributes {
		if *attr.Name == "preferred_username" {
			username = *attr.Value
		}
		if *attr.Name == "email" {
			email = *attr.Value
		}
	}

	return dto.User{
		User:  username,
		Email: email,
	}, nil
}

func NewUserController() {
	UserControllerObj = &UserControllerImpl{}
}
