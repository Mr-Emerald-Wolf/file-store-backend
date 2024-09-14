package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/handlers"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	userRoutes := incomingRoutes.Group("/")
	userRoutes.POST("/register", handlers.CreateUser)
	userRoutes.POST("/login", handlers.LoginUser)
	userRoutes.GET("/user", middleware.VerifyAccessToken(), handlers.GetUser)
}
