-- name: CreatePlayer :one
INSERT INTO players(id, created_at, updated_at, player_name)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2
)
ON CONFLICT (id) DO UPDATE 
SET player_name = EXCLUDED.player_name, updated_at = NOW()
RETURNING *;

-- name: DeletePlayers :exec
DELETE FROM players;

-- name: FetchPlayerByID :one
SELECT * FROM players
WHERE id = $1;