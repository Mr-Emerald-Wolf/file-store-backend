package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/handlers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	userRoutes := incomingRoutes.Group("/")
	userRoutes.POST("/register", handlers.CreateUser)
	userRoutes.POST("/login", handlers.LoginUser)
}
