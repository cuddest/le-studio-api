package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Params describes pagination request.
type Params struct {
	Page   int
	Limit  int
	Offset int
}

// Parse extracts pagination from query params.
func Parse(c *gin.Context) Params {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if p < 1 {
		p = 1
	}
	if l < 1 {
		l = 20
	}
	return Params{Page: p, Limit: l, Offset: (p - 1) * l}
}
