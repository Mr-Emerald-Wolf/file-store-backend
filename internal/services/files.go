package services

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/database"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/db"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/internal/models"
	"github.com/mr-emerald-wolf/21BCE0665_Backend/s3handler"
)

const (
	ChunkSize = 5 * 1024 * 1024 // 5MB per chunk
)

func CreateFileMetaData(newFile models.CreateFileRequest, userID int32) error {

	file := db.CreateFileParams{
		UserID:   userID,
		FileName: newFile.FileName,
		S3Url: pgtype.Text{
			String: newFile.S3Url,
			Valid:  true,
		},
		FileSize: newFile.Size,
		FileType: pgtype.Text{
			String: newFile.FileType,
			Valid:  true,
		},
		IsPublic: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}

	// Create New File
	_, err := database.DB.CreateFile(context.Background(), file)
	if err != nil {
		return err
	}
	return nil
}

func UploadToS3(ctx context.Context, file multipart.File, filename string, userID int32) (string, error) {
	// Generate Filename key
	key := strconv.Itoa(int(userID)) + "/" + filename
	// Initiate multipart upload
	uploadInput := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(key),
	}
	uploadOutput, err := s3handler.S3Client.CreateMultipartUpload(uploadInput)
	if err != nil {
		return "", fmt.Errorf("failed to initiate multipart upload: %v", err)
	}

	var completedParts = make(map[int64]*s3.CompletedPart)
	var partNum int64 = 1
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, 1) // Channel for error handling

	buffer := make([]byte, ChunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			cancelUpload(s3handler.S3Client, uploadOutput)
			return "", fmt.Errorf("failed to read file: %v", err)
		}
		if n == 0 {
			break
		}

		partData := make([]byte, n)
		copy(partData, buffer[:n])

		// Upload the part concurrently
		wg.Add(1)
		go func(partNumber int64, data []byte) {
			defer wg.Done()

			uploadPartInput := &s3.UploadPartInput{
				Bucket:     aws.String(os.Getenv("AWS_BUCKET")),
				Key:        aws.String(key),
				PartNumber: aws.Int64(partNumber),
				UploadId:   uploadOutput.UploadId,
				Body:       bytes.NewReader(data),
			}

			uploadPartOutput, err := s3handler.S3Client.UploadPart(uploadPartInput)
			if err != nil {
				errCh <- fmt.Errorf("failed to upload part %d: %v", partNumber, err)
				cancelUpload(s3handler.S3Client, uploadOutput)
				return
			}

			mu.Lock()
			completedParts[partNumber] = &s3.CompletedPart{
				ETag:       uploadPartOutput.ETag,
				PartNumber: aws.Int64(partNumber),
			}
			mu.Unlock()

		}(partNum, partData)

		partNum++
	}

	// Wait for all parts to upload
	wg.Wait()
	close(errCh)

	// Handle any upload errors
	if len(errCh) > 0 {
		cancelUpload(s3handler.S3Client, uploadOutput)
		return "", fmt.Errorf("multipart upload failed")
	}

	// Convert map to slice and sort by part number
	var sortedParts []*s3.CompletedPart
	for i := int64(1); i < partNum; i++ {
		if part, exists := completedParts[i]; exists {
			sortedParts = append(sortedParts, part)
		}
	}

	// Sort the parts by part number
	sort.Slice(sortedParts, func(i, j int) bool {
		return *sortedParts[i].PartNumber < *sortedParts[j].PartNumber
	})

	// Complete multipart upload
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(os.Getenv("AWS_BUCKET")),
		Key:      aws.String(key),
		UploadId: uploadOutput.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: sortedParts,
		},
	}
	_, err = s3handler.S3Client.CompleteMultipartUpload(completeInput)
	if err != nil {
		cancelUpload(s3handler.S3Client, uploadOutput)
		return "", fmt.Errorf("failed to complete multipart upload: %v", err)
	}

	// Create s3 file url
	s3Url := "https://file-upload-trademarkia-bucket.s3.ap-south-1.amazonaws.com/" + key

	return s3Url, nil
}

func GetFilesByUserID(userID int32) (*[]db.File, error) {

	// Get files by user id
	files, err := database.DB.GetFilesByUserID(context.Background(), userID)

	// Return error if no files exist
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("no files exist for this user")
	} else if err != nil {
		return nil, err
	}

	return &files, nil

}

func cancelUpload(s3Client *s3.S3, uploadOutput *s3.CreateMultipartUploadOutput) {
	_, err := s3Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
		Bucket:   uploadOutput.Bucket,
		Key:      uploadOutput.Key,
		UploadId: uploadOutput.UploadId,
	})
	if err != nil {
		log.Printf("failed to abort upload: %v", err)
	}
}
