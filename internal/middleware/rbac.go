package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "unauthorized", "message": "unauthenticated"}})
			return
		}
		role, _ := roleVal.(string)
		for _, a := range allowed {
			if a == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "forbidden", "message": "forbidden"}})
	}
}
