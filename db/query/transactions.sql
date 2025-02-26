-- name: CreateTransaction :one
INSERT INTO transactions (username, user_email, type, currency, amount, address, tx_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_email, type, currency, amount, status, address, tx_hash, created_at;

-- name: GetTransactionByID :one
SELECT id, username, user_email, type, currency, amount, status, address, tx_hash, created_at
FROM transactions
WHERE id = $1;

-- name: GetTransactionsByUserEmail :many
SELECT id, username, user_email, type, currency, amount, status, address, tx_hash, created_at
FROM transactions
WHERE user_email = $1
ORDER BY created_at DESC;

-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET status = $1
WHERE id = $2;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;
