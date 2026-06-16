-- +goose Up
ALTER TABLE followed_teams
ADD COLUMN tri_code TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE followed_teams
DROP COLUMN tri_code;