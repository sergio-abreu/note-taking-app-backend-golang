package notes

import (
	"context"
	"errors"
)

var (
	ErrEmptyUserID             = errors.New("empty user id")
	ErrInvalidUserIDFormat     = errors.New("invalid user id format")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmptyNoteID             = errors.New("empty note id")
	ErrInvalidNoteIDFormat     = errors.New("invalid note id format")
	ErrNoteNotFound            = errors.New("note not found")
	ErrEmptyReminderID         = errors.New("empty reminder id")
	ErrInvalidReminderIDFormat = errors.New("invalid reminder id format")
	ErrReminderNotFound        = errors.New("reminder not found")
	ErrOnlyOneReminderAllowed  = errors.New("only one reminder allowed for a note")
)

type NotesRepository interface {
	FindUser(ctx context.Context, userID string) (User, error)
	CreateUser(ctx context.Context, user User) error

	FindNote(ctx context.Context, userID, noteID string) (Note, error)
	CreateNote(ctx context.Context, note Note) error
	EditNote(ctx context.Context, note Note) error
	DeleteNote(ctx context.Context, note Note) error
	MarkAsComplete(ctx context.Context, note Note) error
	MarkAsInProgress(ctx context.Context, note Note) error

	FindReminder(ctx context.Context, userID, noteID, reminderID string) (Reminder, error)
	ScheduleReminder(ctx context.Context, reminder Reminder) error
	RescheduleReminder(ctx context.Context, reminder Reminder) error
	DeleteReminder(ctx context.Context, reminder Reminder) error
}
