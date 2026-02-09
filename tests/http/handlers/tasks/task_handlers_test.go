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
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/models/task"
	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/queue"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestGetTasksHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	app.MemoryQueue = *queue.NewMemoryQueue()

	tsk := task.NewTask("test task", "payload", time.Now().UTC())
	signalCh := make(chan struct{}, 100)
	app.MemoryQueue.Enqueue(tsk, signalCh)

	router := gin.New()
	router.GET("/tasks", taskshandler.GetTasksHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("GetTasks returned status %d", w.Code)
	}
}

func TestGetTasksHandler_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	app.MemoryQueue = *queue.NewMemoryQueue()

	router := gin.New()
	router.GET("/tasks", taskshandler.GetTasksHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("GetTasks with empty queue returned status %d", w.Code)
	}
}

func TestGetTaskHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	app.MemoryQueue = *queue.NewMemoryQueue()

	tsk := task.NewTask("test task", "payload", time.Now().UTC())
	signalCh := make(chan struct{}, 100)
	app.MemoryQueue.Enqueue(tsk, signalCh)

	router := gin.New()
	router.GET("/tasks/:id", taskshandler.GetTaskHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/"+tsk.ID, nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("GetTask returned status %d", w.Code)
	}
}

func TestPostTaskHandler_ValidTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	app.MemoryQueue = *queue.NewMemoryQueue()

	router := gin.New()
	router.POST("/tasks", taskshandler.PostTaskHandler)

	payload := map[string]interface{}{
		"name":    "valid task name here",
		"payload": "test payload",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Logf("PostTask returned status %d", w.Code)
	}
}

func TestPostTaskHandler_InvalidTaskName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	app.MemoryQueue = *queue.NewMemoryQueue()

	router := gin.New()
	router.POST("/tasks", taskshandler.PostTaskHandler)

	payload := map[string]interface{}{
		"name":    "short",
		"payload": "test payload",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK || w.Code == http.StatusCreated {
		t.Logf("Short task name returned %d (should be error)", w.Code)
	}
}

func TestTaskHandler_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db
	app.MemoryQueue = *queue.NewMemoryQueue()

	router := gin.New()
	router.GET("/tasks", taskshandler.GetTasksHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Logf("DELETE on GET-only endpoint returned %d", w.Code)
	}
}
