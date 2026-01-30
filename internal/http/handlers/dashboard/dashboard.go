package dashboard

import "github.com/gin-gonic/gin"

func DashboardHandler(c *gin.Context) {
	c.String(200, "hello")
}
