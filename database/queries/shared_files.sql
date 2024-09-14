-- name: CreateSharedFile :one
INSERT INTO shared_files (user_id, file_id, s3_url)
VALUES ($1, $2, $3)
RETURNING id, user_id, file_id, s3_url, shared_at;

-- name: GetSharedFileByID :one
SELECT id, user_id, file_id, s3_url, shared_at
FROM shared_files
WHERE id = $1;

-- name: GetSharedFilesByUserID :many
SELECT id, user_id, file_id, s3_url, shared_at
FROM shared_files
WHERE user_id = $1;

-- name: GetSharedFilesByFileID :many
SELECT id, user_id, file_id, s3_url, shared_at
FROM shared_files
WHERE file_id = $1;

-- name: DeleteSharedFile :exec
DELETE FROM shared_files
WHERE id = $1;

-- name: DeleteSharedFilesByFileID :exec
DELETE FROM shared_files
WHERE file_id = $1;

-- name: DeleteSharedFilesByUserID :exec
DELETE FROM shared_files
WHERE user_id = $1;

-- name: DeleteExpiredSharedFiles :exec
DELETE FROM shared_files
WHERE shared_at < NOW() - INTERVAL '30 minutes';
