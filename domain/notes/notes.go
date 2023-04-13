package notes

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrEmptyTitle = errors.New("empty title")
)

func newNote(title, description string, userID uuid.UUID) (Note, error) {
	if len(title) == 0 {
		return Note{}, ErrEmptyTitle
	}
	return Note{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       title,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}, nil
}

type Note struct {
	ID          uuid.UUID
	Title       string
	Description string
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
