package notes

import (
	"errors"

	"github.com/gofrs/uuid"
)

var (
	ErrNoteDoesntBelongToThisUser     = errors.New("note doesn't belong to this user")
	ErrReminderDoesntBelongToThisUser = errors.New("reminder doesn't belong to this user")
)

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

func (u User) CreateNote(title, description string) (Note, error) {
	return newNote(title, description, u.ID)
}

func (u User) EditNote(note *Note, title, description string) error {
	if err := u.validateNoteBelongsToUser(note); err != nil {
		return err
	}
	return note.edit(title, description)
}

func (u User) MarkNoteAsCompleted(note *Note) error {
	if err := u.validateNoteBelongsToUser(note); err != nil {
		return err
	}
	return note.markAsCompleted()
}

func (u User) MarkNoteAsInProgress(note *Note) error {
	if err := u.validateNoteBelongsToUser(note); err != nil {
		return err
	}
	return note.markAsInProgress()
}

func (u User) CopyNote(note Note) (Note, error) {
	if err := u.validateNoteBelongsToUser(&note); err != nil {
		return Note{}, err
	}
	return note.copy(), nil
}

func (u User) DeleteNote(note Note) error {
	if err := u.validateNoteBelongsToUser(&note); err != nil {
		return err
	}
	return nil
}

func (u User) ScheduleAReminder(note Note, cronExpression, rawEndsAt string, repeats uint) (Reminder, error) {
	if err := u.validateNoteBelongsToUser(&note); err != nil {
		return Reminder{}, err
	}
	return newReminder(note.ID, u.ID, cronExpression, rawEndsAt, repeats)
}

func (u User) RescheduleAReminder(reminder *Reminder, cronExpression, rawEndsAt string, repeats uint) error {
	if err := u.validateReminderBelongsToUser(reminder); err != nil {
		return err
	}
	return reminder.reschedule(cronExpression, rawEndsAt, repeats)
}

func (u User) DeleteReminder(reminder Reminder) error {
	if err := u.validateReminderBelongsToUser(&reminder); err != nil {
		return err
	}
	return nil
}

func (u User) validateNoteBelongsToUser(note *Note) error {
	if u.ID != note.UserID {
		return ErrNoteDoesntBelongToThisUser
	}
	return nil
}

func (u User) validateReminderBelongsToUser(reminder *Reminder) error {
	if u.ID != reminder.UserID {
		return ErrReminderDoesntBelongToThisUser
	}
	return nil
}
