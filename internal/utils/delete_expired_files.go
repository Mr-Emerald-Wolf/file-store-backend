package utils

import (
	"context"
	"log"

	"github.com/mr-emerald-wolf/21BCE0665_Backend/database"
)

func DeleteExpiredSharedFiles() {
	database.DB.DeleteExpiredSharedFiles(context.Background())
	log.Print("Deleting expired shared files")
}

func DeleteExpiredFiles() {
	database.DB.DeleteOldFiles(context.Background())
	log.Print("Deleting expired s3 files")
}
