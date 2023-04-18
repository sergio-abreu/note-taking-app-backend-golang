CREATE TABLE IF NOT EXISTS notes
(
    id          UUID        NOT NULL PRIMARY KEY,
    title       VARCHAR     NOT NULL,
    description VARCHAR,
    completed   BOOLEAN     NOT NULL DEFAULT FALSE,
    user_id     UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);