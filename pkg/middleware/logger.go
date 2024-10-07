package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	end := time.Now()
	latency := end.Sub(start)
	log.Printf("[%s] %s %s %v", c.Request.Method, c.Request.URL.Path, c.ClientIP(), latency)
}
