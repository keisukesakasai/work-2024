package controllers

import (
	"log"
	"net/http"
	"todoapi/app/models"
	"todoapi/app/utils"

	"github.com/gin-gonic/gin"
)

func createTodo(c *gin.Context) {
	utils.Logger(c, "TODO登録")

	var createTodorequest createTodoRequest
	if err := c.BindJSON(&createTodorequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	content := createTodorequest.Content
	user_id := createTodorequest.User_Id
	if err := models.CreateTodo(c, content, user_id); err != nil {
		log.Println(err)
	}

	utils.Logger(c, "TODO登録完了")

	c.JSON(http.StatusOK, gin.H{
		"content": content,
	})
}

func updateTodo(c *gin.Context) {
	utils.Logger(c, "TODO更新")
	var updateTodorequest updateTodoRequest
	if err := c.BindJSON(&updateTodorequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content := updateTodorequest.Content
	user_id := updateTodorequest.User_Id
	todo_id := updateTodorequest.Todo_Id

	if err := models.UpdateTodo(c, content, user_id, todo_id); err != nil {
		log.Println(err)
	}
	utils.Logger(c, "TODO登録完了")

	c.JSON(http.StatusOK, gin.H{
		"content": content,
	})
}

func getTodos(c *gin.Context) {

}

func getTodo(c *gin.Context) {
	utils.Logger(c, "TODO参照")

	var getTodorequest getTodoRequest
	if err := c.BindJSON(&getTodorequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo_id := getTodorequest.Todo_Id
	todo, err := models.GetTodo(c, todo_id)
	if err != nil {
		log.Println(err)
	}
	utils.Logger(c, "TODO参照完了")

	c.JSON(http.StatusOK, gin.H{
		"ID":        todo.ID,
		"Content":   todo.Content,
		"UserID":    todo.UserID,
		"CreatedAt": todo.CreatedAt,
	})
}

func deleteTodo(c *gin.Context) {
	utils.Logger(c, "TODO削除")

	var deleteTodorequest deleteTodoRequest
	if err := c.BindJSON(&deleteTodorequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo_id := deleteTodorequest.Todo_Id
	todo, err := models.GetTodo(c, todo_id)
	if err != nil {
		log.Println(err)
	}
	utils.Logger(c, "TODO参照完了")

	if todo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"resultCode": "Todoがありません",
		})
	} else {
		err := models.DeleteTodo(c, todo_id)
		if err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"resultCode": "ID : " + todo_id + " の Todo を正常に削除しました",
		})
	}
}

func getTodosByUser(c *gin.Context) {
	utils.Logger(c, "ユーザごとのTODO参照")

	var getTodosByUserrequest getTodosByUserRequest
	if err := c.BindJSON(&getTodosByUserrequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id := getTodosByUserrequest.User_Id
	todos, err := models.GetTodosByUser(c, user_id)
	if err != nil {
		log.Println(err)
	}
	utils.Logger(c, "ユーザごとのTODO参照完了")

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}
