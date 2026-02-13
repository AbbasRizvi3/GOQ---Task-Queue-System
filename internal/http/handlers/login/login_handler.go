package login

import (
	"net/http"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/auth"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/utils/authentication"
	"github.com/gin-gonic/gin"
)

func LoginPageHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		data := gin.H{
			"Error":   "",
			"Success": "",
		}
		c.HTML(http.StatusOK, "login.html", data)
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func LoginHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req auth.LoginRequest
		req = auth.LoginRequest{
			Email:    c.PostForm("email"),
			Password: c.PostForm("password"),
		}
		if req.Email == "" || req.Password == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{
				"Error": "Email and password are required",
				"Email": req.Email,
			})
			return
		}
		var storedHash string
		err := app.Databasehandle.QueryRow("SELECT password_hash FROM users WHERE email = $1", req.Email).Scan(&storedHash)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"Error": "invalid email or password",
				"Email": req.Email,
			})
			return
		}
		if !authentication.CheckPasswordHash(req.Password, storedHash) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{
				"Error": "invalid email or password",
				"Email": req.Email,
			})
			return
		}
		err = app.Databasehandle.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&req.ID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Error": "failed to fetch user ID",
				"Email": req.Email,
			})
			return
		}
		token, err := authentication.CreateToken(req.Email, req.ID)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{
				"Error": "failed to create token",
				"Email": req.Email,
			})
			return
		}
		c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/api/dashboard")
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}
