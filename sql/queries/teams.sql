-- name: CreateTeam :one
INSERT INTO teams(id, created_at, updated_at, team_name, tri_code)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;