package logout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", true, true)
	c.SetCookie("session_id", "", -1, "/", "", true, true)

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("Surrogate-Control", "no-store")
	c.Redirect(http.StatusSeeOther, "/login")
}
