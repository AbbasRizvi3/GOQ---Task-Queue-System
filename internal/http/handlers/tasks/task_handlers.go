package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/gin-gonic/gin"
)

func GetTasksHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		rows, err := app.Databasehandle.Query("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks from database"})
			return
		}
		defer rows.Close()

		var tasks []task.RequestTasks
		for rows.Next() {
			var t task.RequestTasks
			if err := rows.Scan(&t.ID, &t.Name, &t.Payload, &t.State, &t.RunAt, &t.CreatedAt, &t.MaxRetries, &t.Retries, &t.NextRunAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan task from database"})
				return
			}
			tasks = append(tasks, t)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while reading tasks from database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tasks": tasks,
		})

	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func GetTaskHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := c.Param("id")
		var t task.RequestTasks
		err := app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Name, &t.Payload, &t.State, &t.RunAt, &t.CreatedAt, &t.MaxRetries, &t.Retries, &t.NextRunAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task from database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"task": t,
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
			t.RunAt = time.Now().UTC()
			fmt.Printf("Set RunAt to now: %v\n", t.RunAt)
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
		err = app.Databasehandle.QueryRow("INSERT INTO tasks (id, name, payload, state, run_at, next_run_at, lease_until, max_retries, retries, error, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
			newTask.ID, newTask.Name, newTask.Payload, newTask.State, newTask.RunAt, newTask.NextRunAt, newTask.LeaseUntil, newTask.MaxRetries, newTask.Retries, newTask.Error, time.Now().UTC(),
		).Scan(&newTask.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert task into database"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Task created successfully",
		})
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func RetryTaskHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		id := c.Param("id")
		var task task.RequestTasks
		err := app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1", id).Scan(&task.ID, &task.Name, &task.Payload, &task.State, &task.RunAt, &task.CreatedAt, &task.MaxRetries, &task.Retries, &task.NextRunAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task from database"})
			return
		}
		if task.State == "leased" || task.State == "completed" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only tasks in failed, dead or canceled state can be retried"})
			return
		} else {
			_, err = app.Databasehandle.Exec("UPDATE tasks SET state = $1, next_run_at = $2, updated_at = $3 WHERE id = $4",
				"ready", time.Now().UTC(), time.Now().UTC(), id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task in database"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Task retried successfully",
			})
		}
	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}

func CancelTaskHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		id := c.Param("id")
		var taskBuffer task.RequestTasks
		err := app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1", id).Scan(&taskBuffer.ID, &taskBuffer.Name, &taskBuffer.Payload, &taskBuffer.State, &taskBuffer.RunAt, &taskBuffer.CreatedAt, &taskBuffer.MaxRetries, &taskBuffer.Retries, &taskBuffer.NextRunAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task from database"})
			return
		}
		if taskBuffer.State == "completed" || taskBuffer.State == "cancelled" || taskBuffer.State == "dead" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only tasks in pending, ready, leased, retry or failed state can be canceled"})
			return
		} else {
			_, err = app.Databasehandle.Exec("UPDATE tasks SET state = $1, updated_at = $2 WHERE id = $3",
				"canceled", time.Now().UTC(), id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task in database"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Task canceled successfully",
			})
		}

	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}
