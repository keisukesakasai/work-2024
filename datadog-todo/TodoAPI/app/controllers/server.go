package controllers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"todoapi/config"

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
	r.POST("/createTodo", createTodo)
	r.POST("/updateTodo", updateTodo)
	r.POST("/deleteTodo", deleteTodo)

	r.POST("/getTodo", getTodo)
	r.POST("/getTodos", getTodos)
	r.POST("/getTodosByUser", getTodosByUser)

	r.Run(":" + serverPort)
}
