package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = os.Getenv("JWT_SECRET")

func verifyToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cannot read claims")
	}

	return &claims, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[len("Bearer "):]
		} else {
			cookieToken, err := c.Cookie("auth_token")
			if err != nil {
				c.HTML(http.StatusMethodNotAllowed, "error.html", gin.H{
					"Code":        http.StatusUnauthorized,
					"Message":     "Unauthorized",
					"Description": "You need to log in to access this page.",
				})
				c.Abort()
				return
			}
			tokenString = cookieToken
		}

		claims, err := verifyToken(tokenString)
		if err != nil {
			c.HTML(http.StatusUnauthorized, "error.html", gin.H{
				"Code":        http.StatusUnauthorized,
				"Message":     "Unauthorized",
				"Description": "Invalid token.",
			})
			c.Abort()
			return
		}
		if username, ok := (*claims)["username"].(string); ok {
			c.Set("username", username)
		}

		var userID string
		if idFloat, ok := (*claims)["id"].(float64); ok {
			userID = fmt.Sprintf("%.0f", idFloat)
		} else if idStr, ok := (*claims)["id"].(string); ok {
			userID = idStr
		}
		if userID != "" {
			c.Set("id", userID)
		}

		c.Next()
	}
}
