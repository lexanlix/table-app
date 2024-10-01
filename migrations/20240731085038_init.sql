-- +goose Up
CREATE TABLE finances
(
    id              UUID NOT NULL PRIMARY KEY,
    main_category   TEXT NOT NULL,
    category        TEXT NOT NULL REFERENCES table_app.category(name) ON UPDATE CASCADE,
    value           INT NOT NULL,
    month           INT NOT NULL,
    year            INT NOT NULL
);

CREATE TABLE category
(
    id              UUID NOT NULL,
    name            TEXT NOT NULL UNIQUE,
    main_category   TEXT NOT NULL,
    priority        INT NOT NULL,

    CONSTRAINT category_pk PRIMARY KEY (main_category, priority)
);

-- +goose Down
DROP TABLE finances;
DROP TABLE category;
