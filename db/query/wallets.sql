-- name: CreateWallet :one
INSERT INTO wallets (username, user_email, currency, balance)
VALUES ($1, $2, $3, $4)
RETURNING id, username, user_email, currency, balance, locked_balance, created_at;

-- name: GetWalletByID :one
SELECT id, username, user_email, currency, balance, locked_balance, created_at
FROM wallets
WHERE id = $1;

-- name: GetWalletsByUserEmail :many
SELECT id, username, user_email, currency, balance, locked_balance, created_at
FROM wallets
WHERE user_email = $1
ORDER BY created_at DESC;

-- name: UpdateWalletBalance :exec
UPDATE wallets
SET balance = $1, locked_balance = $2
WHERE id = $3;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = $1;
