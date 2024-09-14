package cronjobs

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/utils"
)

var scheduler gocron.Scheduler

func RunCronJobs() {
	var err error
	scheduler, err = gocron.NewScheduler(gocron.WithLocation(time.UTC))
	if err != nil {
		log.Fatalf("could not start cron jobs: %v", err)
	}

	// Delete Expired Shared Files every 10 Minutes
	_, err = scheduler.NewJob(
		gocron.CronJob("*/1 * * * *", false),
		gocron.NewTask(utils.DeleteExpiredSharedFiles),
	)
	if err != nil {
		log.Println("Error scheduling shared file deletion:", err)
	}

	// Delete Expired Files stored on s3 twice a day
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(0, 0, 0),
			gocron.NewAtTime(12, 0, 0),
		)),
		gocron.NewTask(utils.DeleteExpiredFiles),
	)
	if err != nil {
		log.Println("Error scheduling daily expired file deletion:", err)
	}

	scheduler.Start()
	log.Println("Cron jobs started successfully")
}

func StopCronJobs() {
	log.Println("Stopping cron jobs")
	if scheduler != nil {
		scheduler.Shutdown()
	}
}
