package logging

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()
		statusCode := c.Writer.Status()
		log.Printf("%s %s - %d", method, path, statusCode)
	}
}
