-- name: FollowPlayer :one
INSERT INTO followed_players(id, created_at, updated_at, user_id, player_id, player_name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: DeletePlayerFollows :exec
DELETE FROM followed_players;

-- name: GetFollowedPlayers :many
SELECT * FROM followed_players
WHERE user_id = $1;

-- name: UnfollowPlayer :exec
DELETE FROM followed_players
WHERE player_id = $1 AND user_id = $2;

