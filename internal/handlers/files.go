package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/db"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/models"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/services"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/utils"
)

func UploadFile(c *gin.Context) {

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

	// Parse file from request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
		return
	}
	defer file.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Upload file in chunks
	s3Url, err := services.UploadToS3(ctx, file, fileHeader.Filename, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HandleError(err))
		return
	}

	// Get File Type
	contentType := fileHeader.Header.Get("Content-Type")

	if contentType == "" {
		contentType = "octet-stream"
	}

	newFile := models.CreateFileRequest{
		FileName: fileHeader.Filename,
		Size:     fileHeader.Size,
		FileType: contentType,
		S3Url:    s3Url,
	}

	// Save file metadata
	err = services.CreateFileMetaData(newFile, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"s3Url":   s3Url,
	})
}

func GetFiles(c *gin.Context) {

	// Parse JWT token
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

	// Get files
	files, err := services.GetFilesByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	// Return file metadata from db
	c.JSON(http.StatusOK, gin.H{
		"message": "all files metadata retrieved successfully",
		"files":   files,
	})
}
