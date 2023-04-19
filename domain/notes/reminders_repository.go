package notes

import (
	"context"
	"errors"
)

var (
	ErrEmptyReminderID         = errors.New("empty reminder id")
	ErrInvalidReminderIDFormat = errors.New("invalid reminder id format")
	ErrReminderNotFound        = errors.New("reminder not found")
)

type RemindersRepository interface {
	FindReminder(ctx context.Context, userID, reminderID string) (Reminder, error)
	ScheduleReminder(ctx context.Context, reminder Reminder) error
	RescheduleReminder(ctx context.Context, reminder Reminder) error
	DeleteReminder(ctx context.Context, reminder Reminder) error
}
