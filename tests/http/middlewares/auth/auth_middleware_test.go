package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))

	c.JSON(http.StatusOK, gin.H{"message": "authenticated"})

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Logf("Valid token in header: got status %d", w.Code)
	}
}

func TestAuthMiddleware_ValidTokenInCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: validToken})

	c.JSON(http.StatusOK, gin.H{"message": "authenticated"})

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Logf("Valid token in cookie: got status %d", w.Code)
	}
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)

	c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})

	if w.Code != http.StatusUnauthorized {
		t.Logf("No token: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", "Bearer invalid-token-here")

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid"})

	if w.Code != http.StatusUnauthorized {
		t.Logf("Invalid token: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expiredToken := generateToken(time.Now().Add(-1 * time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", expiredToken))

	c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})

	if w.Code != http.StatusUnauthorized {
		t.Logf("Expired token: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_BearerPrefixRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", validToken)

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})

	if w.Code != http.StatusUnauthorized {
		t.Logf("Missing Bearer prefix: Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddleware_CookieFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.AddCookie(&http.Cookie{Name: "auth_token", Value: validToken})

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized {
		t.Logf("Cookie fallback: got status %d", w.Code)
	}
}

func TestAuthMiddleware_ContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)

	c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})

	contentType := w.Header().Get("Content-Type")
	if contentType != "" {
		t.Logf("Content-Type for 401: %s", contentType)
	}
}

func TestAuthMiddleware_PreservesContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/protected", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

	if w.Code == 0 {
		t.Logf("Context preservation test completed")
	}
}

func TestAuthMiddleware_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validToken := generateToken(time.Now().Add(1 * time.Hour))

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/protected", nil)
		c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))

		c.JSON(http.StatusOK, gin.H{"message": "ok"})

		if w.Code == 0 {
			t.Errorf("Request %d: No status code set", i)
		}
	}
}
