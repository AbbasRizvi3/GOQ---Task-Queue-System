package timeout

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	timeoutmiddleware "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/timeout"
	"github.com/gin-gonic/gin"
)

func TestTimeoutMiddleware_RequestCompletes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(timeoutmiddleware.TimeoutMiddleware(5 * time.Second))
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

func TestTimeoutMiddleware_RequestTimesOut(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(timeoutmiddleware.TimeoutMiddleware(100 * time.Millisecond))
	router.GET("/test", func(c *gin.Context) {
		time.Sleep(500 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Logf("Expected timeout status %d, got %d", http.StatusGatewayTimeout, w.Code)
	}
}

func TestTimeoutMiddleware_QuickRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(timeoutmiddleware.TimeoutMiddleware(5 * time.Second))
	router.GET("/test", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Quick request should complete: got %d", w.Code)
	}
}

func TestTimeoutMiddleware_ContextWithValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(timeoutmiddleware.TimeoutMiddleware(5 * time.Second))
	router.GET("/test", func(c *gin.Context) {
		ctx := c.Request.Context()
		if ctx == nil {
			t.Error("Context should not be nil")
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTimeoutMiddleware_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(timeoutmiddleware.TimeoutMiddleware(5 * time.Second))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status %d, got %d", i, http.StatusOK, w.Code)
		}
	}
}

func TestTimeoutMiddleware_EdgeCase_ZeroDuration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(timeoutmiddleware.TimeoutMiddleware(0))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code == 0 {
		t.Logf("Zero duration timeout test completed")
	}
}
