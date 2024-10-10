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

		// Extract the token from the "Bearer <token>" format
		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token := parts[1]

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
