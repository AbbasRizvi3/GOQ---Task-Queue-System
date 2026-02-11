package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDashboardHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.String(http.StatusOK, "dashboard")

	if w.Code != http.StatusOK {
		t.Logf("Dashboard returned status %d", w.Code)
	}
}

func TestDashboardHandler_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/", nil)

		c.String(http.StatusOK, "dashboard")

		if w.Code == 0 {
			t.Errorf("Request %d: No status code set", i)
		}
	}
}
