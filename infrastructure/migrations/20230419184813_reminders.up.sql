CREATE TABLE IF NOT EXISTS reminders
(
    id              UUID        NOT NULL PRIMARY KEY,
    note_id         UUID        NOT NULL,
    user_id         UUID        NOT NULL,
    cron_expression VARCHAR     NOT NULL,
    ends_at         TIMESTAMPTZ NULL,
    created_at      TIMESTAMPTZ NOT NULL,
    updated_at      TIMESTAMPTZ NULL,
    FOREIGN KEY (note_id) REFERENCES notes (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE UNIQUE INDEX unique_notes_idx ON reminders(note_id);