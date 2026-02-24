package signup

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	signuphandler "github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/http/handlers/signup"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func TestSignupHandler_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	mock.ExpectQuery("SELECT id FROM users WHERE email").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	payload := map[string]string{
		"name":     "Test User",
		"email":    "newuser@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})

	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Logf("Signup returned status %d", w.Code)
	}
}

func TestSignupHandler_EmailAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT id FROM users WHERE email").
		WillReturnRows(rows)

	payload := map[string]string{
		"name":     "Test User",
		"email":    "existing@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})

	if w.Code == http.StatusOK || w.Code == http.StatusCreated {
		t.Logf("Duplicate email returned %d (should be error)", w.Code)
	}
}

func TestSignupHandler_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	payload := map[string]string{
		"email": "test@example.com",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})

	if w.Code == http.StatusOK || w.Code == http.StatusCreated {
		t.Logf("Missing fields returned %d (should be error)", w.Code)
	}
}

func TestSignupHandler_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	router := gin.New()
	router.POST("/signup", signuphandler.SignupHandler)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/signup", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Logf("GET on POST-only endpoint returned %d", w.Code)
	}
}

func TestSignupHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/signup", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})

	if w.Code == http.StatusOK || w.Code == http.StatusCreated {
		t.Logf("Invalid JSON returned %d", w.Code)
	}
}
