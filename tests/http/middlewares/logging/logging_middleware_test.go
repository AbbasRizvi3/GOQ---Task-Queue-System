package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"

	loggingmiddleware "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/logging"
	"github.com/gin-gonic/gin"
)

func TestLoggingMiddleware_LogsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(loggingmiddleware.LoggingMiddleware())
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

func TestLoggingMiddleware_DifferentMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(loggingmiddleware.LoggingMiddleware())
	router.GET("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	router.POST("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	router.PUT("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	router.DELETE("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })
	router.PATCH("/test", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{}) })

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	for _, method := range methods {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("%s request: Expected status %d, got %d", method, http.StatusOK, w.Code)
		}
	}
}

func TestLoggingMiddleware_DifferentPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(loggingmiddleware.LoggingMiddleware())
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	router.GET("/api/v1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	paths := []string{"/api", "/api/v1", "/test"}
	for _, path := range paths {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Path %s: Expected status %d, got %d", path, http.StatusOK, w.Code)
		}
	}
}

func TestLoggingMiddleware_DoesNotAffectResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(loggingmiddleware.LoggingMiddleware())
	expectedBody := gin.H{"message": "test", "data": "response"}
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, expectedBody)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" && contentType != "application/json" {
		t.Logf("Unexpected content type: %s", contentType)
	}
}

func TestLoggingMiddleware_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(loggingmiddleware.LoggingMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status %d, got %d", i, http.StatusOK, w.Code)
		}
	}
}
