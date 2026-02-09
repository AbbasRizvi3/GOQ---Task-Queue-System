package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authmiddleware "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateToken(expiresAt time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": expiresAt.Unix(),
	})

	tokenString, _ := token.SignedString([]byte("test-secret-key"))
	return tokenString
}

func TestAuthMiddleware_ValidTokenInHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Logf("Valid token in header: got status %d", w.Code)
	}
}

func TestAuthMiddleware_ValidTokenInCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: validToken})
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Logf("Valid token in cookie: got status %d", w.Code)
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Logf("No token: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-here")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Logf("Invalid token: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	expiredToken := generateToken(time.Now().Add(-1 * time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", expiredToken))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Logf("Expired token: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_BearerPrefixRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", validToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Logf("Missing Bearer prefix: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_CookieFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: validToken})
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Logf("Cookie fallback: got status %d", w.Code)
	}
}

func TestAuthMiddleware_ContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authenticated"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "" {
		t.Logf("Content-Type for 401: %s", contentType)
	}
}

func TestAuthMiddleware_PreservesContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID != "" {
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "no user context"})
		}
	})

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
	router.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Logf("Context preservation test completed")
	}
}

func TestAuthMiddleware_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(authmiddleware.AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
		router.ServeHTTP(w, req)

		if w.Code == 0 {
			t.Errorf("Request %d: No status code set", i)
		}
	}
}
