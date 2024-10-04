-- +goose Up
CREATE TABLE last_updated_data
(
    id          TEXT primary key,
    updated_at  TIMESTAMP,
    note        TEXT
);

-- +goose Down
DROP TABLE last_updated_data;