package middleware

import (
	"net/http"
	"strings"

	"cotton-cloud-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Allow unauthenticated access for demo/development
		if authHeader == "" {
			// Check for user_id query param (demo mode)
			if userID := c.Query("user_id"); userID != "" {
				c.Set("userID", userID)
				c.Set("email", "demo@example.com")
				c.Next()
				return
			}

			// No auth provided - use demo user
			c.Set("userID", "demo-user")
			c.Set("email", "demo@example.com")
			c.Next()
			return
		}

		// Validate Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}

// GetUserID extracts user ID from gin context
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return userID.(string)
	}
	return "demo-user"
}

// GetEmail extracts email from gin context
func GetEmail(c *gin.Context) string {
	if email, exists := c.Get("email"); exists {
		return email.(string)
	}
	return "demo@example.com"
}

// RequireAuth strictly requires authentication (no demo mode)
func RequireAuth() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("claims", claims)

		c.Next()
	}
}
