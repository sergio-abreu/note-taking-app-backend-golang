package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewPgCron(db *gorm.DB) PgCron {
	return PgCron{db: db}
}

type PgCron struct {
	db *gorm.DB
}

func (p PgCron) CreateCron(ctx context.Context, reminder notes.Reminder) error {
	cronExpression := reminder.ParseCron()
	endsAt := reminder.ParseEndsAt(cronExpression)
	ends := "NULL"
	if !endsAt.IsZero() {
		ends = "'" + endsAt.Format(time.RFC3339) + "'"
	}
	sql := fmt.Sprintf("SELECT cron.schedule(?, ?, $$ SELECT note_taking_reminder_webhook ('%s', '%s', '%s', %s); $$);", reminder.UserID, reminder.NoteID, reminder.ID, ends)
	return p.db.WithContext(ctx).Exec(sql, reminder.ID, cronExpression).Error
}

func (p PgCron) DeleteCron(ctx context.Context, reminder notes.Reminder) error {
	sql := "SELECT cron.unschedule(?);"
	err := p.db.WithContext(ctx).Exec(sql, reminder.ID).Error
	if pgErr, ok := err.(*pgconn.PgError); ok && strings.Contains(pgErr.Message, "could not find valid entry for job") {
		return nil
	}
	return err
}
