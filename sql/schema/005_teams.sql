-- +goose Up
CREATE TABLE teams (
    id INT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    team_name TEXT NOT NULL,
    tri_code TEXT NOT NULL
);

-- +goose Down
DROP TABLE teams;