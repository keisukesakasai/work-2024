package controllers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"userapi/config"

	"github.com/gin-gonic/gin"
)

var serverPort = config.Config.Port

func StartMainServer() {
	log.Println("info: Start Server" + "port: " + serverPort)

	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// router 設定
	r := gin.New()

	//--- handler 設定
	r.POST("/createUser", createUser)
	r.POST("/getUserByEmail", getUserByEmail)

	r.POST("/encrypt", Encrypt)

	r.Run(":" + serverPort)
}
