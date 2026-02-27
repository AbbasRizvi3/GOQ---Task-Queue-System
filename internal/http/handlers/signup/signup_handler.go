package signup

import (
	"database/sql"
	"net/http"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/auth"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/utils/authentication"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func userExists(email string) (bool, error) {
	var id string
	err := app.Databasehandle.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func insertUser(name, email, hashedPassword string) (string, error) {
	var newUserID string
	err := app.Databasehandle.QueryRow(
		"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		name, email, hashedPassword,
	).Scan(&newUserID)
	return newUserID, err
}

func SignupPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func SignupHandler(c *gin.Context) {
	req := auth.SignupRequest{
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
	exists, err := userExists(req.Email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"Error": "Failed to check existing user",
			"Name":  req.Name,
			"Email": req.Email,
		})
		return
	}
	if exists {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{
			"Error": "Email already exists",
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
	newUserID, err := insertUser(req.Name, req.Email, hashedPassword)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"Error": "Failed to create user",
			"Name":  req.Name,
			"Email": req.Email,
		})
		return
	}
	token, err := authentication.CreateToken(req.Email, newUserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{
			"Error": "Account created, but auto-login failed. Please try logging in manually.",
		})
		return
	}
	c.SetCookie("auth_token", token, 3600*24, "/", "", true, true)
	c.Redirect(http.StatusSeeOther, "/api/dashboard")
}
