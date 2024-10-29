-- +goose Up
CREATE TABLE IF NOT EXISTS chat (
  id BIGSERIAL PRIMARY KEY,
  usernames text[],
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS chat;