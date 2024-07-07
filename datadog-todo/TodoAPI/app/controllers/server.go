package controllers

import (
	"context"
	"log"
	"os"
	"os/signal"
	"todoapi/config"

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
	r.Use(gintrace.Middleware("TodoAPI"))

	//--- handler 設定
	r.POST("/createTodo", createTodo)
	r.POST("/updateTodo", updateTodo)
	r.POST("/deleteTodo", deleteTodo)

	r.POST("/getTodo", getTodo)
	r.POST("/getTodos", getTodos)
	r.POST("/getTodosByUser", getTodosByUser)

	r.Run(":" + serverPort)
}
