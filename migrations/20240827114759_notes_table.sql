-- +goose Up
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    owner TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS notes;
