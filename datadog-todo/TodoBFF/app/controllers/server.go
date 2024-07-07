package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func StartMainServer() {
	log.Println("info: Start Server" + "port: " + serverPort)

	// Datadog Tracer
	tracer.Start()
	defer tracer.Stop()

	// コンテキスト生成
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// router 設定
	r := gin.New()
	r.Use(gintrace.Middleware("TodoBFF"))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// template 設定
	r.LoadHTMLGlob(pathStatic + "/templates/*")
	r.Static("/static/", pathStatic)

	//--- handler 設定
	r.GET("/", top)
	r.GET("/login", getLogin)
	r.POST("/login", postLogin)
	r.GET("/signup", getSignup)
	r.POST("/signup", postSignup)

	rTodos := r.Group("/menu")
	rTodos.Use(checkSession())
	{
		rTodos.GET("/todos", getIndex)
		rTodos.GET("/todos/new", getTodoNew)
		rTodos.POST("/todos/save", postTodoSave)
		rTodos.GET("/todos/edit/:id", parseURL(getTodoEdit))
		rTodos.POST("/todos/update/:id", parseURL(postTodoUpdate))
		rTodos.GET("/todos/delete/:id", parseURL(getTodoDelete))
	}

	r.GET("/logout", getLogout)

	r.Run(":" + serverPort)
}

func checkSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		Logger(c, "セッションチェック開始")

		session := sessions.Default(c)
		LoginInfo.UserID = session.Get("UserId")

		if LoginInfo.UserID == nil {
			Logger(c, LoginInfo.UserID.(string)+" はログインしていません")
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()
		} else {
			Logger(c, LoginInfo.UserID.(string)+" をセッション ID にセットしました")
			c.Set("UserId", LoginInfo.UserID) // ユーザIDをセット
			c.Next()
		}

		Logger(c, "セッションチェック終了")
	}
}
