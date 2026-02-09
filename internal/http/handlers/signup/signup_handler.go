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
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Name == "" || req.Email == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			return
		}
		var existingID int
		err := app.Databasehandle.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		}
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		_, err = app.Databasehandle.Exec("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)", req.Name, req.Email, hashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "signup successful"})
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}
