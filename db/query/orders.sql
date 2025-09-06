-- name: CreateOrder :one
INSERT INTO orders (username, user_email, market_id, type, status, price, amount)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_email, market_id, type, status, price, amount, filled_amount, created_at, updated_at;

-- name: GetOrderByID :one
SELECT id, username, user_email, market_id, type, status, price, amount, filled_amount, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: UpdateOrderStatusAndFilledAmount :exec
UPDATE orders
SET status = $1, filled_amount = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: ListOrders :many
SELECT id, username, user_email, market_id, type, status, price, amount, filled_amount, created_at, updated_at
FROM orders
ORDER BY created_at DESC;
