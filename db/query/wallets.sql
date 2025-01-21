-- name: CreateWallet :one
INSERT INTO wallets (user_id, currency, balance)
VALUES ($1, $2, $3)
RETURNING id, user_id, currency, balance, locked_balance, created_at;

-- name: GetWalletByUserIDAndCurrency :one
SELECT user_id, currency, balance, locked_balance, created_at
FROM wallets
WHERE user_id = $1 AND currency = $2;

-- name: UpdateWalletBalance :exec
UPDATE wallets
SET balance = $1, locked_balance = $2
WHERE user_id = $3 AND currency = $4;

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = $1;
