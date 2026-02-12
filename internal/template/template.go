package template

import "github.com/gin-gonic/gin"

func SetupTemplate(r *gin.Engine) {
	r.LoadHTMLGlob("./internal/static/*.html")
}
