-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: FetchUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
INNER JOIN refresh_token ON user_id = refresh_token.user_id
WHERE refresh_token.token = $1
AND refresh_token.revoked_at IS NULL
AND refresh_token.expires_at > NOW();

-- name: DeleteUsers :exec
DELETE FROM users;