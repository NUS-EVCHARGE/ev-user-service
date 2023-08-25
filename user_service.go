package main

import (
	"flag"
	"github.com/NUS-EVCHARGE/ev-user-service/config"
	"github.com/NUS-EVCHARGE/ev-user-service/controller"
	"github.com/NUS-EVCHARGE/ev-user-service/dao"
	_ "github.com/NUS-EVCHARGE/ev-user-service/docs"
	"github.com/NUS-EVCHARGE/ev-user-service/handler"
	jwt "github.com/akhettar/gin-jwt-cognito"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	r  *gin.Engine
	mw *jwt.AuthMiddleware
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
		err = dao.InitDB(configObj.Dsn)
		if err != nil {
			logrus.WithField("config", configObj).Error("failed to connect to database")
			return
		}
	}
	mw, err = jwt.AuthJWTMiddleware(configObj.Iss, configObj.UserPoolId, configObj.Region)
	if err != nil {
		panic(err)
	}
	controller.NewUserController()
	InitHttpServer(configObj.HttpAddress)
}

func InitHttpServer(httpAddress string) {
	r = gin.Default()
	registerHandler()

	if err := r.Run(httpAddress); err != nil {
		logrus.WithField("error", err).Errorf("http server failed to start")
	}
}

func registerHandler() {
	// use to generate swagger ui
	//	@BasePath	/api/v1
	//	@title		User Service API
	//	@version	1.0
	//	@schemes	http
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// api versioning
	v1 := r.Group("/api/v1")
	// get user info handler
	v1.GET("/user/get_user_info", mw.MiddlewareFunc(), handler.GetUserInfoHandler)
}
