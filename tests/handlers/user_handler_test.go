package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUser(payload models.CreateUserRequest) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *MockService) LoginUser(payload models.LoginUserRequest) (string, error) {
	args := m.Called(payload)
	return args.String(0), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockService := new(MockService)
	mockService.On("CreateUser", mock.Anything).Return(nil)

	// Create a Gin router and register the handler
	r := gin.Default()
	r.POST("/users", func(c *gin.Context) {
		CreateUser(c)
	})

	payload := `{"email": "test@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"message": "user created successfully"}`, w.Body.String())
}

func TestLoginUser(t *testing.T) {
	mockService := new(MockService)
	mockService.On("LoginUser", mock.Anything).Return("mocked-token", nil)

	// Create a Gin router and register the handler
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		LoginUser(c)
	})

	payload := `{"email": "test@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "login successful", "access_token": "mocked-token"}`, w.Body.String())
}
