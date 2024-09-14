package models

type CreateFileRequest struct {
	FileName string
	Size     int64
	FileType string
	S3Url    string
}

type ShareFileRequest struct {
	FileID int32 `uri:"file_id" binding:"required"`
}