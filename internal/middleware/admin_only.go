package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOnly requires admin role claim.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": gin.H{"code": "FORBIDDEN", "message": "Admin role required."}})
			return
		}
		c.Next()
	}
}
