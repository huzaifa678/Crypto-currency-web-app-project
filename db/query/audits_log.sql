-- name: CreateAuditLog :one
INSERT INTO audit_logs (user_email, action, ip_address)
VALUES ($1, $2, $3)
RETURNING id, user_email, action, ip_address, created_at;

-- name: GetAuditLogsByUserEmail :many
SELECT id, user_email, action, ip_address, created_at
FROM audit_logs
WHERE user_email = $1
ORDER BY created_at DESC;

-- name: DeleteAuditLog :exec
DELETE FROM audit_logs
WHERE id = $1;
