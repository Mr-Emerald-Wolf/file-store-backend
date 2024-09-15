package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/db"
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
	
	validator := validator.New()

	if err := validator.Struct(payload); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Please pass in all the required fields", "error": err.Error()})
		return
	}

	err := services.CreateUser(payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
	})
}

func LoginUser(c *gin.Context) {

	var payload models.LoginUserRequest
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	validator := validator.New()

	if err := validator.Struct(payload); err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "Please pass in all the required fields", "error": err.Error()})
		return
	}

	access_token, err := services.LoginUser(payload)

	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.HandleError(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": access_token,
	})
}

func GetUser(c *gin.Context) {

	token, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	user, ok := token.(db.User)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not parse user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "user found successfully",
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
