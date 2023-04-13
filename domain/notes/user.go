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
	if u.ID != note.UserID {
		return ErrNoteDoesntBelongToThisUser
	}
	return note.edit(title, description)
}
