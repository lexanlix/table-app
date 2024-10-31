-- +goose Up
CREATE TABLE note
(
    id          UUID NOT NULL PRIMARY KEY,
    text        TEXT NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE note;