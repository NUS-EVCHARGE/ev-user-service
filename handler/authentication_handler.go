package handler

import (
	"github.com/NUS-EVCHARGE/ev-user-service/controller/authentication"
	"github.com/NUS-EVCHARGE/ev-user-service/controller/encryption"
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	var (
		credentials dto.Credentials
	)

	err := c.BindJSON(&credentials)
	if err != nil {
		logrus.WithField("err", err).Error("error params")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, "Error", "Wrong params"))
		return
	}

	decryptPassword, err := encryption.EncryptionControllerObj.DecryptPassword(credentials.Password)
	if err != nil {
		logrus.WithField("err", err).Error("error decrypting password")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, "Error", "Invalid password"))
		return
	}

	credentials.Password = decryptPassword

	resp, err := authentication.AuthenticationControllerObj.LoginUser(credentials)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			logrus.WithField("err", awsErr.Code()).Error("error authenticating user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, awsErr.Error(), awsErr.Message()))
		} else {
			logrus.WithField("err", err).Error("error authenticating user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, err.Error(), "Please retry again"))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
	return
}

func SignUpHandler(c *gin.Context) {
	var (
		signUpRequest dto.SignUpRequest
	)

	err := c.BindJSON(&signUpRequest)
	if err != nil {
		logrus.WithField("err", err).Error("error params")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, "Error","Wrong params"))
		return
	}

	err = authentication.AuthenticationControllerObj.RegisterUser(signUpRequest)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			logrus.WithField("err", awsErr.Code()).Error("error signing up user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, awsErr.Error(), awsErr.Message()))
		} else {
			logrus.WithField("err", err).Error("error signing up user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, err.Error(), "Please retry again"))
		}
		return
	}
	c.JSON(http.StatusOK, CreateResponse(http.StatusOK, "Success", "Please check your email for confirmation code"))
	return
}

func ConfirmUserHandler(c *gin.Context) {
	var (
		userInfo dto.ConfirmUser
	)

	err := c.BindJSON(&userInfo)
	if err != nil {
		logrus.WithField("err", err).Error("error params")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, "Error","Wrong params"))
		return
	}

	err = authentication.AuthenticationControllerObj.ConfirmUser(userInfo)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			logrus.WithField("err", awsErr.Code()).Error("error confirming user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, awsErr.Error(), awsErr.Message()))
		} else {
			logrus.WithField("err", err).Error("error confirming user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, err.Error(), "Please retry again"))
		}
		return
	}

	c.JSON(http.StatusOK, CreateResponse(http.StatusOK, "Success","User has registered successfully"))
	return
}

func ResendChallengeCodeHandler(c *gin.Context) {
	var (
		resendRequest dto.SignUpResendRequest
	)

	err := c.BindJSON(&resendRequest)
	if err != nil {
		logrus.WithField("err", err).Error("error params")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, "Error","Wrong params"))
		return
	}

	err = authentication.AuthenticationControllerObj.ResendChallengeCode(resendRequest)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			logrus.WithField("err", awsErr.Code()).Error("error resending confirmation")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, awsErr.Error(), awsErr.Message()))
		} else {
			logrus.WithField("err", err).Error("error resending confirmation")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, err.Error(), "Please retry again"))
		}
		return
	}
	c.JSON(http.StatusOK, CreateResponse(http.StatusOK, "Success", "Confirmation code resent"))
	return
}

func LogoutUserHandler(c *gin.Context) {
	// Retrieve the access token from the Authorization header
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		logrus.Error("Authorization header is missing")
		c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, "Error", "Authorization header is missing"))
		return
	}

	// Optional: remove "Bearer " prefix if it exists
	const bearerPrefix = "Bearer "
	if len(accessToken) > len(bearerPrefix) && accessToken[:len(bearerPrefix)] == bearerPrefix {
		accessToken = accessToken[len(bearerPrefix):]
	}

	err := authentication.AuthenticationControllerObj.LogoutUser(accessToken)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			logrus.WithField("err", awsErr.Code()).Error("error logging out user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, awsErr.Error(), awsErr.Message()))
		} else {
			logrus.WithField("err", err).Error("error logging out user")
			c.JSON(http.StatusBadRequest, CreateResponse(http.StatusBadRequest, err.Error(), "Please retry again"))
		}
		return
	}
	c.JSON(http.StatusOK, CreateResponse(http.StatusOK, "", "Logout successful"))
	return
}
