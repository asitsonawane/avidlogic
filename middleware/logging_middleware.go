package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// RequestLogger logs each request made to the server
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture the start time
		startTime := time.Now()

		// Process the request
		c.Next()

		// Calculate the duration
		duration := time.Since(startTime)

		// Log the details of the request
		log.Printf("%s %s | Status: %d | Duration: %v | Client IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}
