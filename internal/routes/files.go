package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/handlers"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/middleware"
)

func FileRoutes(incomingRoutes *gin.Engine) {
	fileRoutes := incomingRoutes.Group("/", middleware.VerifyAccessToken())
	fileRoutes.POST("/upload", handlers.UploadFile)
}
