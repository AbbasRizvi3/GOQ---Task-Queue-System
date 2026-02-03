package routers

import (
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/dashboard"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/healthz"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/login"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/signup"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/static"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/tasks"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/auth"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/logging"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/recovery"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/timeout"
	"github.com/gin-gonic/gin"
	ginprom "github.com/zsais/go-gin-prometheus"
)

const (
	timeoutDuration = 20 * time.Second
)

func SetUpRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(logging.LoggingMiddleware(), timeout.TimeoutMiddleware(timeoutDuration), recovery.RecoveryMiddleware)
	p := ginprom.NewPrometheus("gin")
	p.Use(router)
	router.GET("/", dashboard.DashboardHandler)
	router.GET("/static", static.StaticHandler)
	router.GET("/healthz", healthz.HealthzHandler)
	router.POST("/login", login.LoginHandler)
	router.POST("/signup", signup.SignupHandler)
	group1 := router.Group("/api", auth.AuthMiddleware())
	group1.GET("/tasks", tasks.GetTasksHandler)
	group1.POST("/tasks", tasks.PostTaskHandler)
	group1.GET("/tasks/:id", tasks.GetTaskHandler)
	group1.POST("/tasks/:id/retry", tasks.RetryTaskHandler)
	group1.POST("/tasks/:id/cancel", tasks.CancelTaskHandler)
	return router
}
