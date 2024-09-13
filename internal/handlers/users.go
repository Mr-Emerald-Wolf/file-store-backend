package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/database"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/db"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/models"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var payload models.CreateUserRequest
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	// Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	user := db.CreateUserParams{
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
	}

	// Create New User
	_, err := database.DB.CreateUser(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HandleError(err))
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

	user, err := database.DB.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "password does not match",
		})
		return
	}

	// Generate New Access Token
	access_token, err := utils.CreateToken(user.Email, utils.ACCESS_TOKEN)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": access_token,
	})
}
