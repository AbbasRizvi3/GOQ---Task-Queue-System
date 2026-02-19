package metrics

import (
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime = time.Now()

func MetricsHandler(c *gin.Context) {
	if c.Query("format") == "json" {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		c.JSON(http.StatusOK, gin.H{
			"goroutines":   runtime.NumGoroutine(),
			"memory_alloc": m.Alloc / 1024 / 1024,
			"memory_total": m.TotalAlloc / 1024 / 1024,
			"memory_sys":   m.Sys / 1024 / 1024,
			"num_gc":       m.NumGC,
			"cpus":         runtime.NumCPU(),
			"uptime":       time.Since(startTime).String(),
		})
		return
	}

	accept := c.GetHeader("Accept")
	if strings.Contains(accept, "text/html") {
		c.HTML(http.StatusOK, "metrics.html", nil)
		return
	}

	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
