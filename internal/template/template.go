package template

import "github.com/gin-gonic/gin"

func SetupTemplate(r *gin.Engine) {
	r.LoadHTMLGlob("/home/dev/Desktop/go-proj/GOQ---Task-Queue-System/internal/static/*.html")
}
