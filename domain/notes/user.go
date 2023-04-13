package notes

import (
	"errors"

	"github.com/gofrs/uuid"
)

var (
	ErrNoteDoesntBelongToThisUser = errors.New("note doesn't belong to this user")
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

func (u User) validateNoteBelongsToUser(note *Note) error {
	if u.ID != note.UserID {
		return ErrNoteDoesntBelongToThisUser
	}
	return nil
}
