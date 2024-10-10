package middleware

import (
	"net/http"
	"strings"

	"avidlogic/auth" // Import the auth package for JWT validation

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks the JWT token for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// If the Bearer prefix is missing, add it automatically
		if !strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = "Bearer " + authHeader
		}

		// Extract the token
		token := strings.Split(authHeader, "Bearer ")[1]

		// Validate the token
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user ID from the token in the context for further use
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
