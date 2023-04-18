CREATE TABLE IF NOT EXISTS users
(
    id         UUID        NOT NULL PRIMARY KEY,
    name       VARCHAR     NOT NULL,
    email      VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NULL
);