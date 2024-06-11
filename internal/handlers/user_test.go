package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sudhir512kj/ecommerce_backend/config"
	"github.com/sudhir512kj/ecommerce_backend/internal/models"
	"github.com/sudhir512kj/ecommerce_backend/internal/repository"
)

type postgresDatabase struct {
	Db *sql.DB
}

var once sync.Once
var conf *config.Config

var db *postgresDatabase

func TestUserHandler_Register(t *testing.T) {
	// Set up the test environment
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.SSLMode,
			conf.Db.TimeZone,
		)

		dbI, err := sql.Open("postgres", dsn)
		if err != nil {
			panic("failed to connect database")
		}

		db = &postgresDatabase{Db: dbI}
	})
	// require.NoError(t, err)
	defer db.Db.Close()

	userRepo := repository.NewUserRepository(db.Db)
	h := NewUserHandler(userRepo)

	t.Run("should register a new user", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		r := gin.Default()
		r.POST("/users", h.Register)

		body, _ := json.Marshal(map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john@example.com",
			"password":   "password123",
		})
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var resp models.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "John", resp.FirstName)
		assert.Equal(t, "Doe", resp.LastName)
		assert.Equal(t, "john@example.com", resp.Email)
		assert.Contains(t, resp.Permissions, "seller")
	})

	t.Run("should return error when email is already registered", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		r := gin.Default()
		r.POST("/users", h.Register)

		body, _ := json.Marshal(map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john@example.com",
			"password":   "password123",
		})
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "email already registered")
	})
}

func TestUserHandler_Login(t *testing.T) {
	// Set up the test environment
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.SSLMode,
			conf.Db.TimeZone,
		)

		dbI, err := sql.Open("postgres", dsn)
		if err != nil {
			panic("failed to connect database")
		}

		db = &postgresDatabase{Db: dbI}
	})
	defer db.Db.Close()

	userRepo := repository.NewUserRepository(db.Db)
	h := NewUserHandler(userRepo)

	// Register a new user
	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/users", h.Register)

	body, _ := json.Marshal(map[string]interface{}{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john@example.com",
		"password":   "password123",
	})
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	// Test the Login handler
	t.Run("should login a user", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		r := gin.Default()
		r.POST("/login", h.Login)

		body, _ := json.Marshal(map[string]interface{}{
			"email":    "john@example.com",
			"password": "password123",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("should return error when email or password is invalid", func(t *testing.T) {
		// Arrange
		w := httptest.NewRecorder()
		r := gin.Default()
		r.POST("/login", h.Login)

		body, _ := json.Marshal(map[string]interface{}{
			"email":    "john@example.com",
			"password": "wrongpassword",
		})
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email or password")
	})
}

// Add more test cases for other user handler functions
