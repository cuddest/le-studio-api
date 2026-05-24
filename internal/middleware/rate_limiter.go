package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter is a placeholder limiter middleware.
func RateLimiter(_ int, _ time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) { c.Next() }
}
