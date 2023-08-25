package controller

import (
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/golang-jwt/jwt"
)

type UserController interface {
	GetUserInfo(jwtToken *jwt.Token) (dto.User, error)
}

type UserControllerImpl struct {
}

var (
	UserControllerObj UserController
)

func (u *UserControllerImpl) GetUserInfo(token *jwt.Token) (dto.User, error) {
	// Cast the claims
	claims := token.Claims.(jwt.MapClaims)

	return dto.User{
		User:  claims["cognito:username"].(string),
		Email: claims["email"].(string),
	}, nil
}

func NewUserController() {
	UserControllerObj = &UserControllerImpl{}
}
