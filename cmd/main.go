package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/config"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/database"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.CheckEnv()
	cfg := config.LoadConfig()
	database.InitDB(cfg.DatabaseConfig)
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start the server
	r.Run(":8080")
}
