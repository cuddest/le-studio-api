package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS applies configured CORS policy.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowAny := len(allowedOrigins) == 0
	allowed := map[string]struct{}{}
	for _, origin := range allowedOrigins {
		trimmed := strings.TrimSpace(origin)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			allowAny = true
		}
		allowed[trimmed] = struct{}{}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// If configured to allow any origin, echo the request Origin when present
		// so browsers can accept credentialed requests. Otherwise only allow
		// configured origins.
		if origin != "" {
			if allowAny {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				if _, ok := allowed[origin]; ok {
					c.Header("Access-Control-Allow-Origin", origin)
				}
			}
			// Allow cookies/credentials when an explicit origin is returned
			c.Header("Access-Control-Allow-Credentials", "true")
			// Instruct caches/proxies that response varies by Origin
			c.Header("Vary", "Origin")
		} else if allowAny {
			// No Origin header (e.g. same-origin requests) — allow all
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// Allow common headers used by the admin UI
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With, Access-Control-Request-Headers")
		c.Header("Access-Control-Expose-Headers", "Content-Range, X-Total-Count")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")

		// Handle preflight
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
