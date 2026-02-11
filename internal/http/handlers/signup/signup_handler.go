package signup

import (
	"net/http"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SignupPageHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "signup.html", nil)
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func SignupHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req auth.SignupRequest
		req = auth.SignupRequest{
			Name:     c.PostForm("name"),
			Email:    c.PostForm("email"),
			Password: c.PostForm("password"),
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"Error": "All fields are required",
				"Name":  req.Name,
				"Email": req.Email,
			})
			return
		}
		var existingID int
		err := app.Databasehandle.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
		if err == nil {
			c.HTML(http.StatusBadRequest, "signup.html", gin.H{
				"Error": "Email already exists",
				"Name":  req.Name,
				"Email": req.Email,
			})
			return
		}
		if err != nil && err.Error() != "sql: no rows in result set" {
			c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
				"Error": "Failed to check existing user",
				"Name":  req.Name,
				"Email": req.Email,
			})
			return
		}
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
				"Error": "Failed to hash password",
				"Name":  req.Name,
				"Email": req.Email,
			})
			return
		}
		_, err = app.Databasehandle.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", req.Name, req.Email, hashedPassword)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
				"Error": "Failed to create user",
				"Name":  req.Name,
				"Email": req.Email,
			})
			return
		}
		c.HTML(http.StatusOK, "signup.html", gin.H{
			"Success": "Account created successfully! Redirecting to login...",
			"Name":    "",
			"Email":   "",
		})
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}
