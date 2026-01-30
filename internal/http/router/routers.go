package routers

import (
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/dashboard"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/healthz"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/static"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/tasks"
	"github.com/gin-gonic/gin"
	ginprom "github.com/zsais/go-gin-prometheus"
)

func SetUpRoutes() *gin.Engine {
	router := gin.Default()
	p := ginprom.NewPrometheus("gin")
	p.Use(router)
	router.GET("/", dashboard.DashboardHandler)
	router.GET("/static", static.StaticHandler)
	router.GET("/healthz", healthz.HealthzHandler)
	group1 := router.Group("/api")
	group1.GET("/tasks", tasks.GetTasksHandler)
	group1.POST("/tasks", tasks.PostTaskHandler)
	group1.GET("/tasks/:id", tasks.GetTaskHandler)
	group1.POST("/tasks/:id/retry", tasks.RetryTaskHandler)
	group1.POST("/tasks/:id/cancel", tasks.CancelTaskHandler)
	return router
}
