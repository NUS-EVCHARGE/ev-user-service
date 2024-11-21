package handler

import (
	"fmt"
	"github.com/NUS-EVCHARGE/ev-user-service/controller/authentication"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func GetAccessToken(c *gin.Context) string {
	// Retrieve the access token from the Authorization header
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		logrus.Error("Authorization header is missing")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest,"Authorization header is missing"))
		c.Abort()
		return ""
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(accessToken, bearerPrefix) {
		logrus.Error("Invalid authorization header format")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		c.Abort()
		return ""
	}

	// Optional: remove "Bearer " prefix if it exists
	if len(accessToken) > len(bearerPrefix) && accessToken[:len(bearerPrefix)] == bearerPrefix {
		accessToken = accessToken[len(bearerPrefix):]
	}
	return accessToken
}

func AuthMiddlewareHandler(c *gin.Context) {
	// Retrieve the access token from the Authorization header
	accessToken := GetAccessToken(c)

	userInfo, err := authentication.AuthenticationControllerObj.GetUserInfo(accessToken)
	if err != nil {
		logrus.WithField("err", err).Error("error validating token")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest,fmt.Sprintf("%v", err)))
		c.Abort()
		return
	}

	for _, us := range userInfo.UserAttributes {
		if *us.Name == "email" {
			//user, err := provider.ProviderControllerObj.GetProvider(*us.Value)
			if err != nil {
				c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest,fmt.Sprintf("%v", err)))
				c.Abort()
				return
			}
			//c.Set("provider", provider)
			return
		}
	}
	c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest,"no user details found"))
	c.Abort()
}
