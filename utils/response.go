package utils

import (
	"github.com/gin-gonic/gin"
)

// JSONSuccess writes a success response with data
func JSONSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{"status": "success", "data": data})
}

// JSONMessage writes a success response with message
func JSONMessage(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"status": "success", "message": message})
}

// JSONError writes an error response
func JSONError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"status": "error", "message": message})
}
