package cache

import "github.com/gin-gonic/gin"

func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("Surrogate-Control", "no-store")

		c.Next()
	}
}
