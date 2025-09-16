-- +goose Up
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);

-- +goose Down
DROP TABLE posts;
