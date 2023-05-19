CREATE TABLE IF NOT EXISTS notes
(
    id          UUID        NOT NULL PRIMARY KEY,
    title       VARCHAR     NOT NULL,
    description VARCHAR,
    completed   BOOLEAN     NOT NULL DEFAULT FALSE,
    user_id     UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX created_at_idx ON notes (user_id, created_at);