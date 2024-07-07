package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func top(c *gin.Context) {
	Logger(c, "TOP画面取得")
	generateHTML(c, "hello", "top", "layout", "top", "public_navbar", "footer")
}

func getIndex(c *gin.Context) {
	Logger(c, "TODO画面取得")

	UserId, isExist := c.Get("UserId")
	if !isExist {
		log.Println("セッションが存在していません")
	}

	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	Logger(c, "UserAPI /getUserByEmail にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req1, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/getUserByEmail",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	rsp1, err := cli.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := io.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI getTodosByUser への Post
	user_id := strconv.Itoa(responseGetUser.ID)
	jsonStr2 := `{"user_id":"` + string(user_id) + `"}`

	Logger(c, "TodoAPI /getTodosByEmail にポスト")
	cli = httptrace.WrapClient(http.DefaultClient)
	req2, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpTodoAPI+"/getTodosByUser",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	rsp2, err := cli.Do(req2)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = io.ReadAll(rsp2.Body)
	var getTodosByUserresponse getTodosByUserResponse
	err = json.Unmarshal(byteArr, &getTodosByUserresponse)
	if err != nil {
		log.Println(err)
	}

	var user User
	user.Name = responseGetUser.Name
	user.Todos = getTodosByUserresponse.Todos

	Logger(c, "TODO画面取得")
	generateHTML(c, user, "index", "layout", "private_navbar", "index", "footer")
}

func getTodoNew(c *gin.Context) {
	Logger(c, "TODO作成画面取得")
	generateHTML(c, nil, "todoNew", "layout", "private_navbar", "todo_new", "footer")
}

func postTodoSave(c *gin.Context) {
	Logger(c, "TODO保存")

	UserId, isExist := c.Get("UserId")
	if !isExist {
		log.Println("セッションが存在していません")
	}

	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	Logger(c, "UserAPI /getUserByEmail にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req1, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/getUserByEmail",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	rsp1, err := cli.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI createTodo への Post
	user_id := strconv.Itoa(responseGetUser.ID)
	content := c.Request.PostFormValue("content")

	Logger(c, "TodoAPI /createTodo にポスト")
	jsonStr2 := `{"Content":"` + content + `",
	"User_Id":"` + user_id + `"}`

	cli = httptrace.WrapClient(http.DefaultClient)
	req2, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpTodoAPI+"/createTodo",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	rsp2, err := cli.Do(req2)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var getTodosByUserresponse getTodosByUserResponse
	err = json.Unmarshal(byteArr, &getTodosByUserresponse)
	if err != nil {
		log.Println(err)
	}

	Logger(c, "TODO画面にリダイレクト")
	c.Redirect(http.StatusFound, "/menu/todos")
}

func getTodoEdit(c *gin.Context, id int) {
	Logger(c, "TODO編集画面取得")

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	UserId, _ := c.Get("UserId")
	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	Logger(c, "UserAPI /getUserByEmail にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req1, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/getUserByEmail",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	rsp1, err := cli.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI getTodo への Post
	todo_id := strconv.Itoa(id)
	jsonStr2 := `{"todo_id":"` + todo_id + `"}`

	Logger(c, "TodoAPI /getTodo にポスト")
	cli = httptrace.WrapClient(http.DefaultClient)
	req2, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpTodoAPI+"/getTodo",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	rsp2, err := cli.Do(req2)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var getTodoresponse getTodoResponse
	err = json.Unmarshal(byteArr, &getTodoresponse)
	if err != nil {
		log.Println(err)
	}

	Logger(c, "TODO編集画面取得")
	generateHTML(c, getTodoresponse, "todoEdit", "layout", "private_navbar", "todo_edit", "footer")
}

func postTodoUpdate(c *gin.Context, id int) {
	Logger(c, "TODO更新")

	err := c.Request.ParseForm()
	if err != nil {
		log.Println(err)
	}

	UserId, _ := c.Get("UserId")
	//--- UserAPI getUserByEmail への Post
	email := UserId.(string)
	jsonStr1 := `{"Email":"` + email + `"}`

	Logger(c, "UserAPI /getUserByEmail にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req1, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpUserApi+"/getUserByEmail",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	rsp1, err := cli.Do(req1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var responseGetUser ResponseGetUser
	err = json.Unmarshal(byteArr, &responseGetUser)
	if err != nil {
		log.Println(err)
	}

	//--- TodoAPI updateTodo への Post
	content := c.Request.PostFormValue("content")
	user_id := strconv.Itoa(responseGetUser.ID)
	todo_id := strconv.Itoa(id)
	jsonStr2 := `{"Content":"` + content + `",
	"User_Id":"` + user_id + `",
	"Todo_Id":"` + todo_id + `"}`

	Logger(c, "TodoAPI /updateTodo にポスト")

	cli = httptrace.WrapClient(http.DefaultClient)
	req2, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpTodoAPI+"/updateTodo",
		bytes.NewBuffer([]byte(jsonStr2)),
	)
	rsp2, err := cli.Do(req2)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp2.Body.Close()

	byteArr, _ = ioutil.ReadAll(rsp2.Body)
	var updateTodoresponse updateTodoResponse
	err = json.Unmarshal(byteArr, &updateTodoresponse)
	if err != nil {
		log.Println(err)
	}

	Logger(c, "TODO画面にリダイレクト")
	c.Redirect(http.StatusFound, "/menu/todos")
}

func getTodoDelete(c *gin.Context, id int) {
	Logger(c, "TODO削除")

	//--- TodoAPI deleteTodo への Post
	todo_id := strconv.Itoa(id)
	jsonStr1 := `{"todo_id":"` + todo_id + `"}`

	Logger(c, "TodoAPI /deleteTodo にポスト")
	cli := httptrace.WrapClient(http.DefaultClient)
	req1, _ := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		EpTodoAPI+"/deleteTodo",
		bytes.NewBuffer([]byte(jsonStr1)),
	)
	rsp1, err := cli.Do(req1)
	if err != nil {
		log.Println(err)
		return
	}
	defer rsp1.Body.Close()

	byteArr, _ := ioutil.ReadAll(rsp1.Body)
	var deleteTodoresponse deleteTodoResponse
	err = json.Unmarshal(byteArr, &deleteTodoresponse)
	if err != nil {
		log.Println(err)
	}

	Logger(c, "TODO画面にリダイレクト")
	c.Redirect(http.StatusFound, "/menu/todos")
}
