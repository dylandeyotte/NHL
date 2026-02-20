-- name: CreateRefreshToken :one
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;

-- name: RevokeToken :one
UPDATE refresh_token
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING *;