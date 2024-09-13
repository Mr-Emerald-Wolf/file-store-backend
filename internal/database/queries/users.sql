-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET username = $1, email = $2, password_hash = $3, updated_at = NOW()
WHERE id = $4
RETURNING id, username, email, password_hash, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

