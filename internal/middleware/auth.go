package middleware

import (
	"net/http"
	"strings"

	"libratrack/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		token := ""
		if strings.HasPrefix(raw, "Bearer ") {
			token = strings.TrimPrefix(raw, "Bearer ")
		}
		if token == "" {
			if ck, err := c.Cookie("libratrack_token"); err == nil {
				token = ck
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "missing token"}})
			return
		}
		claims, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "invalid token"}})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
