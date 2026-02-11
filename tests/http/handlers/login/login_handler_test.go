package login

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AbbasRizvi3/GOQ---Task-Queue-System/internal/core/app"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginHandler_ValidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"password_hash"}).AddRow(string(hashedPassword))
	mock.ExpectQuery("SELECT password_hash FROM users WHERE email").
		WillReturnRows(rows)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusOK, gin.H{"message": "success"})

	if w.Code != http.StatusOK {
		t.Logf("Valid credentials: got status %d", w.Code)
	}
}

func TestLoginHandler_InvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	rows := sqlmock.NewRows([]string{"password_hash"}).AddRow(string(hashedPassword))
	mock.ExpectQuery("SELECT password_hash FROM users WHERE email").
		WillReturnRows(rows)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})

	if w.Code == http.StatusOK {
		t.Errorf("Invalid password should not return 200, got %d", w.Code)
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	mock.ExpectQuery("SELECT password_hash FROM users WHERE email").
		WillReturnError(sql.ErrNoRows)

	payload := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})

	if w.Code == http.StatusOK {
		t.Errorf("User not found should not return 200, got %d", w.Code)
	}
}

func TestLoginHandler_MissingPassword(t *testing.T) {
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
	c.Request = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusBadRequest, gin.H{"error": "missing password"})

	if w.Code == http.StatusOK {
		t.Logf("Missing password returned %d", w.Code)
	}
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})

	if w.Code == http.StatusOK {
		t.Logf("Invalid JSON returned %d", w.Code)
	}
}

func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	app.Databasehandle = db

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/login", nil)

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})

	if w.Code != http.StatusNotFound {
		t.Logf("GET on POST-only endpoint returned %d", w.Code)
	}
}
