-- name: GetGoogleUserByEmail :one
SELECT * FROM google_auth WHERE email = $1 LIMIT 1;

-- name: GetGoogleUserByProviderID :one
SELECT * FROM google_auth WHERE provider_id = $1 LIMIT 1;

-- name: CreateGoogleUser :one
INSERT INTO google_auth (
    email, username, provider, provider_id, role
) VALUES (
    $1, $2, 'google', $3, 'user'
)
RETURNING *;