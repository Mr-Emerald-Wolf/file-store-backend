package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/config"
	cronjobs "github.com/mr-emerald-wolf/21BCE0665_Backend/cron_jobs"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/database"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/routes"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/s3handler"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.CheckEnv()
	cfg := config.LoadConfig()
	database.InitDB(cfg.DatabaseConfig)
	database.NewRepository(cfg.RedisConfig)
	s3handler.InitializeS3Session(cfg.AWSConfig)
}

func main() {
	app := gin.Default()

	// Register Routes
	routes.UserRoutes(app)
	routes.FileRoutes(app)

	// Start CronJobs
	cronjobs.RunCronJobs()
	defer cronjobs.StopCronJobs()

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start the server
	app.Run(":8080")
}
