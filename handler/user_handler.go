package handler

import (
	"github.com/NUS-EVCHARGE/ev-user-service/controller"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"net/http"
)

//	@Summary		Health Check
//	@Description 	perform health check status
//	@Tags 			Health Check
//	@Accept 		json
//	@Produce 		json
//	@Success 		200	{object}	map[string]interface{}	"returns a welcome message"
func GetHealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, CreateResponse("Welcome to ev-user-service", "ok"))
	return
}

//	@Summary		Get User Info from AWS cognito
//	@Description	get user info by authentication header
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.User	"returns a user object"
//	@Router			/user/get_user_info [get]
//	@Param			authentication	header	string	yes	"jwtToken of the user"
func GetUserInfoHandler(c *gin.Context) {
	tokenStr, _ := c.Get("JWT_TOKEN")
	user, err := controller.UserControllerObj.GetUserInfo(tokenStr.(*jwt.Token))
	if err != nil {
		// todo: change to common library
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func CreateResponse(message, body string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"body":    body,
	}
}
