package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	taskshandler "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/tasks"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestGetTasksHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	columns := []string{"id", "name", "payload", "state", "run_at", "created_at", "max_retries", "retries", "next_run_at"}

	rows := sqlmock.NewRows(columns).AddRow(
		"task-123", "test task", "test payload", "pending", time.Now(), time.Now(), 3, 0, nil,
	)

	mock.ExpectQuery("^SELECT (.+) FROM tasks WHERE user_id = \\$1").WillReturnRows(rows)

	router := gin.New()
	router.GET("/tasks", func(c *gin.Context) {
		c.Set("id", "test-user-id")
		taskshandler.GetTasksHandler(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTasks returned status %d. Body: %s", w.Code, w.Body.String())
	}
}
func TestGetTasksHandler_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	columns := []string{
		"id", "user_id", "name", "payload", "state", "run_at",
		"next_run_at", "lease_until", "max_retries", "retries",
		"error", "created_at", "updated_at",
	}
	rows := sqlmock.NewRows(columns)

	mock.ExpectQuery("^SELECT (.+) FROM tasks").WillReturnRows(rows)

	router := gin.New()
	router.GET("/tasks", func(c *gin.Context) {
		c.Set("id", "test-user-id")
		taskshandler.GetTasksHandler(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetTasks with empty queue returned status %d. Body: %s", w.Code, w.Body.String())
	}
}
func TestPostTaskHandler_ValidTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	rows := sqlmock.NewRows([]string{"id"}).AddRow("new-id-1")
	mock.ExpectQuery("^INSERT INTO tasks").WillReturnRows(rows)

	payload := map[string]interface{}{
		"name":    "valid task name here",
		"payload": "test payload",
	}
	body, _ := json.Marshal(payload)

	router := gin.New()
	router.POST("/tasks", func(c *gin.Context) {
		c.Set("id", "test-user-id")
		taskshandler.PostTaskHandler(c)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Errorf("PostTask returned status %d. Error body: %s", w.Code, w.Body.String())
	}
}
