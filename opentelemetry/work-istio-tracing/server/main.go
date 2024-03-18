package main

import (
	"fmt"
	"math/rand"
	"net/http"
	logging "server/internal/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	r := gin.New()
	r.Use(
		otelgin.Middleware("queryFruit"),
	)
	r.GET("/", Handler)
	r.Run(":8080")
}

func Handler(c *gin.Context) {
	randomWord := generateRandomWord(c)
	c.String(http.StatusOK, randomWord)
}

func generateRandomWord(c *gin.Context) string {
	ctx := c.Request.Context()
	logger := logging.GetLoggerFromCtx(ctx)

	// print Header
	headersStr := ""
	for key, values := range c.Request.Header {
		for _, value := range values {
			headersStr += fmt.Sprintf("%s: %s; ", key, value)
		}
	}
	headersStr = strings.TrimRight(headersStr, "; ")
	logger.Infof("Received headers @ Go Server(queryFruit): %s", headersStr)

	words := []string{"apple", "banana", "cherry", "date", "elderberry"}

	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))

	word := words[rand.Intn(len(words))]
	logger.Infof("response: %s", word)
	return word
}
