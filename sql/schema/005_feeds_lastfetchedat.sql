-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
DROP TABLE feeds DROP COLUMN last_fetched_at;