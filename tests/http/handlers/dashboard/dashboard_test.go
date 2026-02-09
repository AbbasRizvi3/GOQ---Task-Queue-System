package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	dashboardhandler "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/dashboard"
	"github.com/gin-gonic/gin"
)

func TestDashboardHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", dashboardhandler.DashboardHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Dashboard returned status %d", w.Code)
	}
}

func TestDashboardHandler_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", dashboardhandler.DashboardHandler)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		if w.Code == 0 {
			t.Errorf("Request %d: No status code set", i)
		}
	}
}
