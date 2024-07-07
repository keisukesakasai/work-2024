package utils

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(c *gin.Context, msg string) {
	start := time.Now()
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()
	logger.Info("Logger",
		zap.Int("status", c.Writer.Status()),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()),
		zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		zap.Duration("elapsed", time.Since(start)),
		zap.String("message", msg),
	)
}
