package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/models"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/services"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/utils"
)

func CreateUser(c *gin.Context) {
	var payload models.CreateUserRequest
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	// Use the service directly or via dependency injection
	err := services.CreateUser(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

func LoginUser(c *gin.Context) {
	var payload models.LoginUserRequest
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	accessToken, err := services.LoginUser(payload)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": accessToken,
	})
}
