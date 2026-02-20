-- +goose Up
CREATE TABLE followed_players (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    player_id INT NOT NULL,
    player_name TEXT NOT NULL,
    UNIQUE (user_id, player_id),
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY (player_id)
    REFERENCES players(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE followed_players;