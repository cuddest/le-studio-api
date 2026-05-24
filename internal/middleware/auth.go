package middleware

import (
	"net/http"
	"strings"
	jwtpkg "le-studio-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// Auth validates bearer token and injects claims.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		a := c.GetHeader("Authorization")
		parts := strings.SplitN(a, " ", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "Missing or invalid token."}})
			return
		}
		claims, err := jwtpkg.ParseAccessToken(secret, parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": gin.H{"code": "UNAUTHORIZED", "message": "Missing or invalid token."}})
			return
		}
		c.Set("userID", claims.Subject)
		c.Set("role", claims.Role)
		c.Next()
	}
}
