-- +goose Up
CREATE TABLE finances
(
    id              SERIAL NOT NULL PRIMARY KEY,
    main_category    TEXT NOT NULL,
    category        TEXT NOT NULL,
    value           INT NOT NULL,
    month           INT NOT NULL,
    year            INT NOT NULL
);

CREATE TABLE category
(
    id              SERIAL NOT NULL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    main_category   TEXT NOT NULL
)

-- +goose Down
DROP TABLE finances;
DROP TABLE category;
