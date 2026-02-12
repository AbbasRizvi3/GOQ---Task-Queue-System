package landing

import "github.com/gin-gonic/gin"

func LandingPageHandler(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.HTML(405, "error.html", gin.H{
			"Code":        405,
			"Message":     "Method Not Allowed",
			"Description": "The method is not allowed for the requested URL.",
		})
		return
	}
	c.HTML(200, "landing.html", nil)
}
