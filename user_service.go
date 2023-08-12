package main

import (
	"ev-user-service/config"
	"ev-user-service/controller"
	"ev-user-service/dao"
	"ev-user-service/util"
	"flag"
	"fmt"
	jwt "github.com/akhettar/gin-jwt-cognito"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	r   *gin.Engine
	err error
	mw  *jwt.AuthMiddleware
)

func main() {
	var (
		configFile string
	)
	flag.StringVar(&configFile, "config", "config.yaml", "configuration file of this service")
	flag.Parse()

	// init configurations
	configObj, err := config.ParseConfig(configFile)
	if err != nil {
		logrus.WithField("error", err).WithField("filename", configFile).Error("failed to init configurations")
		return
	}

	// init db
	if false {
		err = dao.InitDb(configObj.Dsn)
		if err != nil {
			logrus.WithField("config", configObj).Error("failed to connect to database")
			return
		}
	}

	userController = controller.NewUserControl()

	InitHttpServer(configObj.HttpAddress)
}

func InitHttpServer(httpAddress, iss, userPoolId, region string) {
	r = gin.Default()
	mw, err = jwt.AuthJWTMiddleware(iss, userPoolId, region)
	if err != nil {
		panic(err)
	}
	r.POST("/register", RegisterService))

	// proxy request
	r.POST("/booking/create", mw.MiddlewareFunc(), func(context *gin.Context) {

	})
	if err := r.Run(httpAddress); err != nil {
		logrus.WithField("error", err).Errorf("http server failed to start")
	}
}

func RegisterService(c *gin.Context) {
	var (
		url, _                 = c.GetQuery("url")
		serviceName, _         = c.GetQuery("service_name")
		command, _             = c.GetQuery("command")
		method, _              = c.GetQuery("method")
		serviceUrl, gatewayUrl = generateUrl(url, serviceName, command)
	)
	switch method {
	// booking/create-booking
	case "POST":
		r.POST(gatewayUrl, mw.MiddlewareFunc(), func(context *gin.Context) {
			var (
				body map[string]interface{}
				query map[string]string
			)

			err := c.BindJSON(&body)
			if err != nil {
				c.JSON(http.StatusBadRequest, CreateResponse("error binding json", fmt.Sprintf("%v", err)))
				return
			}

			c.BindQuery(&query)
			util.LaunchHttpRequest("POST", serviceUrl, )
		})
	}
}

func generateUrl(url, serviceName, command string) (string, string) {
	return fmt.Sprintf("%v/%v/%v", url, serviceName, command), fmt.Sprintf("%v/%v", serviceName, command)
}

func CreateResponse(message, body string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"body":    body,
	}
}
