package recovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	recoverymiddleware "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/recovery"
	"github.com/gin-gonic/gin"
)

func TestRecoveryMiddleware_NormalRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(recoverymiddleware.RecoveryMiddleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRecoveryMiddleware_CatchesPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(recoverymiddleware.RecoveryMiddleware)
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestRecoveryMiddleware_MultipleHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(recoverymiddleware.RecoveryMiddleware)
	router.GET("/safe", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/panic", func(c *gin.Context) {
		panic("test")
	})
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/safe", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Safe handler: Expected status %d, got %d", http.StatusOK, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Panic handler: Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
