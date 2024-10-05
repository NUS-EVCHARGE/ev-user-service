package user

import (
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	NewUserController()
}

func TestGetUserInfo(t *testing.T) {
	setup()
	var (
		username = "sweiyang"
		email = "e0014576@u.nus.edu"
	)
	mockToken := &jwt.Token{
		Raw:       "",
		Method:    nil,
		Header:    nil,
		Claims:    jwt.MapClaims{
			"cognito:username": username,
			"email": email,
		},
		Signature: "",
		Valid:     false,
	}
	user,err := UserControllerObj.GetUserInfo(mockToken)

	assert.Nil(t, err)
	assert.Equal(t, user, dto.User{
		User:  username,
		Email: email,
	})
}