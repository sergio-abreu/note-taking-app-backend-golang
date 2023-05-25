CREATE TABLE IF NOT EXISTS reminders
(
    id           UUID        NOT NULL PRIMARY KEY,
    note_id      UUID        NOT NULL,
    user_id      UUID        NOT NULL,
    start_date   DATE        NOT NULL,
    start_time   VARCHAR     NOT NULL,
    timezone     VARCHAR     NOT NULL,
    interval     VARCHAR     NOT NULL,
    week_days    VARCHAR     NULL,
    ends_after_n int         NOT NULL,
    ends_at      TIMESTAMPTZ NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (note_id) REFERENCES notes (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX unique_notes_idx ON reminders (note_id);