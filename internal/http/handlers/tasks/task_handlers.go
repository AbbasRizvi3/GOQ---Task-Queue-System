package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) string {
	id, exists := c.Get("id")
	if !exists {
		return ""
	}
	return fmt.Sprintf("%v", id)
}

func GetTasksHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := c.MustGet("id")
		rows, err := app.Databasehandle.Query("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE user_id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks from database"})
			return
		}
		defer func() {
			if err := rows.Close(); err != nil {
				fmt.Printf("Error closing rows: %v\n", err)
			}
		}()

		var tasks []task.HandleTask
		for rows.Next() {
			var t task.HandleTask
			var nextRunAtNull sql.NullTime
			if err := rows.Scan(&t.ID, &t.Name, &t.Payload, &t.State, &t.RunAt, &t.CreatedAt, &t.MaxRetries, &t.Retries, &nextRunAtNull); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan task from database"})
				return
			}
			if nextRunAtNull.Valid {
				t.NextRunAt = nextRunAtNull.Time
			}
			if t.Payload == "" || t.Payload == "null" {
				t.Payload = "No payload"
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
		user_id := c.MustGet("id")
		var t task.HandleTask
		var nextRunAtNull sql.NullTime
		err := app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1 AND user_id = $2", id, user_id).Scan(&t.ID, &t.Name, &t.Payload, &t.State, &t.RunAt, &t.CreatedAt, &t.MaxRetries, &t.Retries, &nextRunAtNull)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task from database"})
			return
		}
		if nextRunAtNull.Valid {
			t.NextRunAt = nextRunAtNull.Time
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
		id := c.MustGet("id")
		var t task.CreateTaskRequest
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
		var nextRunAtDB interface{} = nil
		if !newTask.NextRunAt.IsZero() {
			nextRunAtDB = newTask.NextRunAt
		}
		err = app.Databasehandle.QueryRow("INSERT INTO tasks (id, name, payload, state, run_at, next_run_at, lease_until, max_retries, retries, error, updated_at, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id",
			newTask.ID, newTask.Name, newTask.Payload, newTask.State, newTask.RunAt, nextRunAtDB, newTask.LeaseUntil, newTask.MaxRetries, newTask.Retries, newTask.Error, time.Now(), id,
		).Scan(&newTask.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert task into database"})
			return
		}
		userID := getUserID(c)
		app.WebsocketChannelManager.BroadcastJSON(userID, gin.H{
			"type":    "TASK_CREATED",
			"payload": newTask,
		})

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
		user_id := c.MustGet("id")
		var taskk task.HandleTask
		var nextRunAtNull sql.NullTime
		err := app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1 AND user_id = $2", id, user_id).Scan(&taskk.ID, &taskk.Name, &taskk.Payload, &taskk.State, &taskk.RunAt, &taskk.CreatedAt, &taskk.MaxRetries, &taskk.Retries, &nextRunAtNull)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task from database"})
			return
		}
		if nextRunAtNull.Valid {
			taskk.NextRunAt = nextRunAtNull.Time
		}
		if taskk.State == "leased" || taskk.State == "completed" || taskk.State == "dead" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only tasks in failed, or canceled state with <3 retries can be retried"})
			return
		} else {
			if taskk.Retries == taskk.MaxRetries {
				c.JSON(http.StatusBadRequest, gin.H{"error": "task has exceeded max retries"})
				return
			}
			_, err = app.Databasehandle.Exec("UPDATE tasks SET state = $1, next_run_at = $2, updated_at = $3, retries = $4 WHERE id = $5 AND user_id = $6",
				"pending", time.Now(), time.Now(), taskk.Retries+1, id, user_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task in database"})
				return
			}
			userID := getUserID(c)
			var updatedTask task.HandleTask
			var nextRunAtNull sql.NullTime
			err = app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1 AND user_id = $2", id, userID).
				Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.Payload, &updatedTask.State, &updatedTask.RunAt, &updatedTask.CreatedAt, &updatedTask.MaxRetries, &updatedTask.Retries, &nextRunAtNull)

			if err != nil {
				fmt.Printf("Error fetching updated task for broadcast: %v\n", err)
			} else {
				if nextRunAtNull.Valid {
					updatedTask.NextRunAt = nextRunAtNull.Time
				}
				app.WebsocketChannelManager.BroadcastJSON(userID, gin.H{
					"type":    "TASK_UPDATED",
					"payload": updatedTask,
				})
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
		user_id := c.MustGet("id")
		var taskBuffer task.HandleTask
		var nextRunAtNull sql.NullTime
		err := app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1 AND user_id = $2", id, user_id).Scan(&taskBuffer.ID, &taskBuffer.Name, &taskBuffer.Payload, &taskBuffer.State, &taskBuffer.RunAt, &taskBuffer.CreatedAt, &taskBuffer.MaxRetries, &taskBuffer.Retries, &nextRunAtNull)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch task from database"})
			return
		}
		if taskBuffer.State == "completed" || taskBuffer.State == "canceled" || taskBuffer.State == "dead" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only tasks in pending, leased, retry or failed state can be canceled"})
			return
		} else {
			_, err = app.Databasehandle.Exec("UPDATE tasks SET state = $1, next_run_at = NULL, updated_at = $2 WHERE id = $3 AND user_id = $4",
				"canceled", time.Now(), id, user_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task in database"})
				return
			}
			select {
			case app.CancelSignal <- id:
				fmt.Printf("Sent cancel signal for task %s to worker\n", id)
			default:
			}

			userID := getUserID(c)
			var updatedTask task.HandleTask
			var nextRunAtNull sql.NullTime
			err = app.Databasehandle.QueryRow("SELECT id, name, payload, state, run_at, created_at, max_retries, retries, next_run_at FROM tasks WHERE id = $1 AND user_id = $2", id, userID).
				Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.Payload, &updatedTask.State, &updatedTask.RunAt, &updatedTask.CreatedAt, &updatedTask.MaxRetries, &updatedTask.Retries, &nextRunAtNull)

			if err == nil {
				if nextRunAtNull.Valid {
					updatedTask.NextRunAt = nextRunAtNull.Time
				}
				app.WebsocketChannelManager.BroadcastJSON(userID, gin.H{
					"type":    "TASK_UPDATED",
					"payload": updatedTask,
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Task canceled successfully",
			})
		}

	} else {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
	}
}
