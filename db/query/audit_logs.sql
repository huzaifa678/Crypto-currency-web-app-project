-- name: CreateAuditLog :one
INSERT INTO audit_logs (username, user_email, action, ip_address)
VALUES ($1, $2, $3, $4)
RETURNING id, username, user_email, action, ip_address, created_at;

-- name: GetAuditLogsByUserEmail :many
SELECT id, username, user_email, action, ip_address, created_at
FROM audit_logs
WHERE user_email = $1
ORDER BY created_at DESC;

-- name: DeleteAuditLog :exec
DELETE FROM audit_logs
WHERE id = $1;
