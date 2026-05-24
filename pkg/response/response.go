package response

import "github.com/gin-gonic/gin"

// Meta defines pagination metadata.
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// OK returns success response.
func OK(c *gin.Context, data any) { c.JSON(200, gin.H{"success": true, "data": data}) }

// Created returns created success response.
func Created(c *gin.Context, data any) { c.JSON(201, gin.H{"success": true, "data": data}) }

// Paginated returns paginated success response.
func Paginated(c *gin.Context, data any, meta Meta) {
	c.JSON(200, gin.H{"success": true, "data": data, "meta": meta})
}

// Error returns standard error response.
func Error(c *gin.Context, status int, code, message string, fields any) {
	c.JSON(status, gin.H{"success": false, "error": gin.H{"code": code, "message": message, "fields": fields}})
}
