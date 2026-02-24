package healthz

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HealthzHandler(c *gin.Context) {
	accept := c.GetHeader("Accept")

	if strings.Contains(accept, "application/json") {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		return
	}

	c.HTML(http.StatusOK, "healthz.html", gin.H{
		"Status":  "Operational",
		"Message": "All systems are running smoothly. The task queue is active and processing requests.",
	})
}
