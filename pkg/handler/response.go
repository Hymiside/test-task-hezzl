package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// responseWithError отдает ответ с ошибкой и JSON
func responseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{
		"status":      "error",
		"description": message,
	})
}

// responseSuccessful отдает 200 ответ с JSON
func responseSuccessful(c *gin.Context, message interface{}) {

	fmt.Println(message, "123")

	c.AbortWithStatusJSON(200, gin.H{
		"status":      "success",
		"description": message,
	})
}