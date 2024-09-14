package models

type CreateFileRequest struct {
	FileName string
	Size     int64
	FileType string
	S3Url    string
}
