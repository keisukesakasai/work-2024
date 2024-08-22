package main

import (
	logging "datadog-log-sample-app/internal/log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.GET("/", Handler)
	r.Run(":80")
}

func Handler(c *gin.Context) {
	logger := logging.NewLogger()

	logger.Debug("Reserved request.")
	randomWord := generateRandomWord(c)

	c.String(http.StatusOK, randomWord)
}

func generateRandomWord(c *gin.Context) string {
	logger := logging.NewLogger()

	words := []string{"apple", "banana", "cherry", "coconut", "strawberry"}
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	word := words[rand.Intn(len(words))]
	logger.Info("Response is " + word)

	logger.Error("Error detected !!")
	return word
}
