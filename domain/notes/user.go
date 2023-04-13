package notes

import "github.com/gofrs/uuid"

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

func (u User) CreateNote(title, description string) (Note, error) {
	return newNote(title, description, u.ID)
}
