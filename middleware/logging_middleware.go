package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs details about each incoming request
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time for the request
		startTime := time.Now()

		// Process the request
		c.Next()

		// Duration of the request
		duration := time.Since(startTime)

		// Log the request details
		log.Printf("[Request] %s %s | Status: %d | Duration: %v | IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}
