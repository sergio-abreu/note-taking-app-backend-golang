package notes

import (
	"context"
	"errors"
)

var (
	ErrInvalidReminderID = errors.New("invalid reminder id")
	ErrReminderNotFound  = errors.New("reminder not found")
)

type RemindersRepository interface {
	FindReminder(ctx context.Context, userID, reminderID string) (Reminder, error)
	ScheduleReminder(ctx context.Context, reminder Reminder) error
	RescheduleReminder(ctx context.Context, reminder Reminder) error
	DeleteReminder(ctx context.Context, reminder Reminder) error
}
