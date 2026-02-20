-- +goose Up
CREATE TABLE followed_teams(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    team_name TEXT NOT NULL,
    user_id UUID UNIQUE NOT NULL,
    team_id INT NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY (team_id)
    REFERENCES teams(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE followed_teams;