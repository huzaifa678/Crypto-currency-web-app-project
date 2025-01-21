-- name: CreateTrade :one
INSERT INTO trades (buy_order_id, sell_order_id, market_id, price, amount, fee)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, buy_order_id, sell_order_id, market_id, price, amount, fee, created_at;

-- name: GetTradeByID :one
SELECT id, buy_order_id, sell_order_id, market_id, price, amount, fee, created_at
FROM trades
WHERE id = $1;

-- name: GetTradesByMarketID :many
SELECT id, buy_order_id, sell_order_id, market_id, price, amount, fee, created_at
FROM trades
WHERE market_id = $1
ORDER BY created_at DESC;

-- name: DeleteTrade :exec
DELETE FROM trades
WHERE id = $1;
