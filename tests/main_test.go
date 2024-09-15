package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zeebo/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestHomepageHandler(t *testing.T) {
	mockResponse := `{"message":"pong"}`
	r := SetUpRouter()
	r.GET("localhost:8080/")
	req, _ := http.NewRequest("GET", "ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}
