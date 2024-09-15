package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	existingFile, err := services.CheckIfFileExists(fileHeader.Filename, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HandleError(err))
		return
	}
	if existingFile != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "File already exists",
			"s3Url":   existingFile.S3Url,
		})
		return
	}
	// Upload file in chunks
	s3Url, err := services.UploadToS3(ctx, file, fileHeader.Filename, user.Uuid)
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

func ShareFile(c *gin.Context) {
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

	var payload models.ShareFileRequest
	if err := c.ShouldBindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	publicUrl, err := services.ShareFile(payload.FileID, user.ID, user.Uuid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "file shared successfully",
		"public_url": publicUrl,
	})
}

func SearchFile(c *gin.Context) {
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

	fileName := c.Query("file_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	fileType := c.Query("file_type")

	// Parse Date if Given
	var parsedStartDate, parsedEndDate *time.Time
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
			return
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			return
		}
		parsedStartDate = &start
		parsedEndDate = &end
	}

	var files *[]db.File
	var err error

	switch {
	case fileName != "":
		files, err = services.SearchFileByName(fileName, user.ID)
	case parsedStartDate != nil && parsedEndDate != nil:
		files, err = services.SearchFileByDate(*parsedStartDate, *parsedEndDate, user.ID)
	case fileType != "":
		files, err = services.SearchFileByType(fileType, user.ID)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter (file_name, start_date/end_date, file_type) must be provided"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "files found successfully",
		"files":   files,
	})
}

func UpdateFile(c *gin.Context) {

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

	var file_id models.ShareFileRequest
	if err := c.ShouldBindUri(&file_id); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	var payload models.UpdateFileName
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	if payload.FileName == "" {
		c.JSON(http.StatusBadRequest, utils.HandleError(fmt.Errorf("filename cannot be empty")))
		return
	}

	err := services.UpdateFileMetaData(payload.FileName, file_id.FileID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File updated successfully",
		"New File": payload.FileName,
	})
}

func DeleteFile(c *gin.Context) {

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

	var file_id models.DeleteFileRequest
	if err := c.ShouldBindUri(&file_id); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	err := services.DeleteFileMetaData(file_id.FileID, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}
