-- name: CreateUser :one
INSERT INTO users (email, password_hash, role, is_verified)
VALUES ($1, $2, $3, $4)
RETURNING id, email, created_at, updated_at, role, is_verified;

-- name: GetUserByID :one
SELECT id, email, password_hash, created_at, updated_at, role, is_verified
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at, role, is_verified
FROM users
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET password_hash = $1, updated_at = CURRENT_TIMESTAMP, is_verified = $2
WHERE id = $3;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
