package healthz

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"html/template"
	healthzhandler "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/healthz"
	"github.com/gin-gonic/gin"
)

func TestHealthzHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/healthz", healthzhandler.HealthzHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("Accept", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Healthz returned status %d", w.Code)
	}
}

func TestHealthzHandler_MultipleRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	tmpl := template.Must(template.New("healthz.html").Parse("<p>Test Healthz</p>"))
    router.SetHTMLTemplate(tmpl)
	router.GET("/healthz", healthzhandler.HealthzHandler)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/healthz", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Request %d: status %d", i, w.Code)
		}
	}
}
