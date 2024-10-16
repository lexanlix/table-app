-- +goose Up
CREATE TABLE account
(
    id          UUID NOT NULL PRIMARY KEY,
    name        TEXT NOT NULL,
    sum         INT NOT NULL,
    note        TEXT,
    is_in_sum   BOOL NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE account;