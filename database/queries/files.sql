-- name: CreateFile :one
INSERT INTO files (user_id, file_name, s3_url, file_size, file_type, is_public)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFileByID :one
SELECT id, user_id, file_name, s3_url, file_size, file_type, upload_date, last_accessed, is_public
FROM files
WHERE id = $1;

-- name: GetFileByName :one
SELECT id, user_id, file_name, s3_url, file_size, file_type, upload_date, last_accessed, is_public
FROM files
WHERE file_name = $1;

-- name: GetFilesByUserID :many
SELECT id, user_id, file_name, s3_url, file_size, file_type, upload_date, last_accessed, is_public
FROM files
WHERE user_id = $1
ORDER BY upload_date DESC;

-- name: UpdateFile :one
UPDATE files
SET file_name = $1, last_accessed = NOW()
WHERE id = $2 AND user_id = $3
RETURNING *;

-- name: SearchFilesByName :many
SELECT * FROM files
WHERE user_id = $1 AND file_name ILIKE '%' || $2 || '%';

-- name: SearchFilesByDate :many
SELECT * FROM files
WHERE user_id = $1 AND upload_date BETWEEN $2 AND $3;

-- name: SearchFilesByType :many
SELECT * FROM files
WHERE user_id = $1 AND file_type ILIKE '%' || $2 || '%';

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1 AND user_id = $2;

-- name: DeleteOldFiles :exec
DELETE FROM files
WHERE created_at < NOW() - INTERVAL '30 days';
