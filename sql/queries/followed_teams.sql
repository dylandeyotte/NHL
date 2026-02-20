-- name: FollowTeam :one
INSERT INTO followed_teams(id, created_at, updated_at, team_name, user_id, team_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: DeleteTeamFollows :exec
DELETE FROM followed_teams;

-- name: FetchTeamByTriCode :one
SELECT * FROM teams 
WHERE tri_code = $1;

-- name: UnfollowTeam :exec
DELETE FROM followed_teams
WHERE user_id = $1 AND team_id = $2;

-- name: GetFollowedTeam :one
SELECT * FROM followed_teams
WHERE user_id = $1;