package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func getSignup(c *gin.Context) {
	Logger(c, "ユーザ登録画面取得")
	generateHTML(c, nil, "signup", "layout", "signup", "public_navbar", "footer")
}

func postSignup(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//--- UserAPI createUser への Post
	name := c.Request.PostFormValue("name")
	email := c.Request.PostFormValue("email")
	password := c.Request.PostFormValue("password")

	jsonStr := `{"Name":"` + name + `",
	"Email":"` + email + `",
	"PassWord":"` + password + `"}`

	Logger(c, "UserAPI /createUser にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/createUser",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	rsp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp.Body.Close()
	body, _ := ioutil.ReadAll(rsp.Body)
	log.Println(string(body))

	UserId := email
	login(c, UserId)

	Logger(c, "TODO画面にリダイレクト")
	c.Redirect(http.StatusMovedPermanently, "/menu/todos")
}

func getLogin(c *gin.Context) {
	Logger(c, "ログイン画面取得")
	generateHTML(c, nil, "login", "layout", "login", "public_navbar", "footer")
}

func postLogin(c *gin.Context) {
	Logger(c, "ログイン")
	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//--- UserAPI getUserByEmail への Post
	email := c.Request.PostFormValue("email")
	jsonStr := `{"Email":"` + email + `"}`

	Logger(c, "UserAPI /getUserByEmail にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/getUserByEmail",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	rsp, err := cli.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- UserAPI encrypt への Post
	password := c.Request.PostFormValue("password")
	jsonStr = `{"PassWord":"` + password + `"}`

	Logger(c, "UserAPI /encrypt にポスト")
	cli = httptrace.WrapClient(http.DefaultClient)
	req, _ = http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/encrypt",
		bytes.NewBuffer([]byte(jsonStr)),
	)
	rsp, err = cli.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp.Body)
	var responseEncrypt ResponseEncrypt
	err = json.Unmarshal(byteArr, &responseEncrypt)
	if err != nil {
		log.Println(err)
	}

	if responseGetUser.ID == 0 {
		log.Println("ユーザがいません")

		Logger(c, "ログイン画面にリダイレクト")
		c.Redirect(http.StatusFound, "/login")
	} else if responseEncrypt.PassWord == responseGetUser.PassWord {
		UserId := c.PostForm("email")
		login(c, UserId)

		Logger(c, "TODO画面にリダイレクト")
		c.Redirect(http.StatusMovedPermanently, "/menu/todos")
	} else {
		log.Println("PW が間違っています")

		Logger(c, "ログイン画面にリダイレクト")
		c.Redirect(http.StatusFound, "/login")
	}
}

func getLogout(c *gin.Context) {
	Logger(c, "ログアウト")
	logout(c)

	Logger(c, "TOP画面にリダイレクト")
	c.Redirect(http.StatusMovedPermanently, "/")
}

func login(c *gin.Context, UserId string) {
	Logger(c, "ログイン処理...")

	session := sessions.Default(c)

	Logger(c, "セッション設定")
	session.Set("UserId", UserId)

	Logger(c, "セッション保存")
	session.Save()

	Logger(c, "ログイン完了")
}

func logout(c *gin.Context) {
	Logger(c, "ログアウト処理...")

	session := sessions.Default(c)

	Logger(c, "セッションクリア")

	session.Clear()

	Logger(c, "セッション保存")
	session.Save()

	Logger(c, "ログアウト完了")
}
