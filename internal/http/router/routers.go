package routers

import (
	"net/http"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/dashboard"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/healthz"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/login"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/logout"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/signup"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/tasks"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/auth"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/middlewares/cache"
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
	router.Use(logging.LoggingMiddleware(), timeout.TimeoutMiddleware(timeoutDuration), recovery.RecoveryMiddleware, cache.NoCacheMiddleware())
	p := ginprom.NewPrometheus("gin")
	p.Use(router)
	router.GET("/healthz", healthz.HealthzHandler)
	router.GET("/login", login.LoginPageHandler)
	router.POST("/login", login.LoginHandler)
	router.GET("/signup", signup.SignupPageHandler)
	router.POST("/signup", signup.SignupHandler)
	router.GET("/logout", logout.LogoutHandler)
	group2 := router.Group("/", auth.AuthMiddleware(), cache.NoCacheMiddleware())
	group2.GET("/tasks/:id", tasks.GetTaskDetailHandler)

	group1 := router.Group("/api", auth.AuthMiddleware(), cache.NoCacheMiddleware())
	group1.GET("/dashboard", dashboard.DashboardHandler)
	group1.GET("/tasks", tasks.GetTasksHandler)
	group1.POST("/tasks", tasks.PostTaskHandler)
	group1.GET("/tasks/:id", tasks.GetTaskHandler)
	group1.POST("/tasks/:id/retry", tasks.RetryTaskHandler)
	group1.POST("/tasks/:id/cancel", tasks.CancelTaskHandler)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "error.html", gin.H{
			"Code":        http.StatusNotFound,
			"Message":     "Page Not Found",
			"Description": "The page you're looking for doesn't exist.",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.HTML(405, "error.html", gin.H{
			"Code":        http.StatusMethodNotAllowed,
			"Message":     "Method Not Allowed",
			"Description": "The request method is not allowed for this resource.",
		})
	})
	return router
}
