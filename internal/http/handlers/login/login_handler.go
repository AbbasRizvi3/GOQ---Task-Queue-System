package login

import (
	"net/http"
	"os"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = os.Getenv("JWT_SECRET")
var tokenExpiryHours = time.Hour * 24

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(tokenExpiryHours).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginPageHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}

func LoginHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req auth.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var storedHash string
		err := app.Databasehandle.QueryRow("SELECT password_hash FROM users WHERE email = $1", req.Email).Scan(&storedHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		if !checkPasswordHash(req.Password, storedHash) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		token, err := createToken(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}
		c.SetCookie("auth_token", token, 3600*24, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "login successful"})
	} else {
		c.String(http.StatusMethodNotAllowed, "method not allowed")
	}
}
