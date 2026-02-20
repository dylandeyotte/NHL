-- +goose Up
CREATE TABLE players (
    id INT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    player_name TEXT NOT NULL
);

-- +goose Down
DROP TABLE players;