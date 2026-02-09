package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}
