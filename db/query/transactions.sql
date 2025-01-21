-- name: CreateTransaction :one
INSERT INTO transactions (user_id, type, currency, amount, address, tx_hash)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, type, currency, amount, status, address, tx_hash, created_at;

-- name: GetTransactionByID :one
SELECT id, user_id, type, currency, amount, status, address, tx_hash, created_at
FROM transactions
WHERE id = $1;

-- name: GetTransactionsByUserID :many
SELECT id, user_id, type, currency, amount, status, address, tx_hash, created_at
FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateTransactionStatus :exec
UPDATE transactions
SET status = $1
WHERE id = $2;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;
