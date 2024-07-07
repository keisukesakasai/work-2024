package controllers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"userapi/config"

	"github.com/gin-gonic/gin"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var serverPort = config.Config.Port

func StartMainServer() {
	log.Println("info: Start Server" + "port: " + serverPort)

	// Datadog Tracer
	tracer.Start()
	defer tracer.Stop()

	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// router 設定
	r := gin.New()
	r.Use(gintrace.Middleware("UserAPI"))

	//--- handler 設定
	r.POST("/createUser", createUser)
	r.POST("/getUserByEmail", getUserByEmail)

	r.POST("/encrypt", Encrypt)

	r.Run(":" + serverPort)
}
