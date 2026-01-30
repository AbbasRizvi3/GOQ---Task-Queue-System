package tasks

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/gin-gonic/gin"
)

// get tasks for now reads the results queue only, later it will be changed to sql db making it the single source of truth

func GetTasksHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.JSON(http.StatusOK, gin.H{
			"tasks": app.ResultQueue.GetAllTasks(),
		})
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func GetTaskHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"task": app.ResultQueue.GetTaskByID(id),
		})
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func PostTaskHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		var t task.RequestTask
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if t.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "task name is required"})
			return
		}
		if t.RunAt.IsZero() {
			t.RunAt = time.Now()
		}
		if len(t.Name) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "task name must be at least 8 characters long"})
			return
		}
		payloadBytes, err := json.Marshal(t.Payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal payload"})
			return
		}
		newTask := task.NewTask(t.Name, string(payloadBytes), t.RunAt)
		app.MemoryQueue.Enqueue(newTask, app.SignalCh)
		c.JSON(http.StatusOK, gin.H{
			"message": "Task created successfully",
		})
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func RetryTaskHandler(c *gin.Context) {
}

func CancelTaskHandler(c *gin.Context) {
}
