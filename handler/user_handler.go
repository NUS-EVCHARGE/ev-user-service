package handler

import (
	"github.com/NUS-EVCHARGE/ev-user-service/controller/user"
	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, CreateResponse(
		http.StatusOK,"Welcome to ev-user-service", "success"))
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
	tokenStr := GetAccessToken(c)
	userInfo, err := user.UserControllerObj.GetUserInfo(tokenStr)
	if err != nil {
		// todo: change to common library
		logrus.WithField("err", err).Error("error getting user")
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, userInfo)
	return
}


func CreateResponse(statusCode int, data interface{}, message ...string) map[string]interface{} {
	response := map[string]interface{}{
		"status":  statusCode,
		"data":    data,
	}
	if len(message) > 0 {
		response["message"] = message[0]
	}
	return response
}