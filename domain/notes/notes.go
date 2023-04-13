package notes

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrEmptyTitle             = errors.New("empty title")
	ErrNotIsAlreadyCompleted  = errors.New("note is already completed")
	ErrNotIsAlreadyInProgress = errors.New("note is already in progress")
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
	Completed   bool
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (n *Note) edit(title, description string) error {
	if len(title) == 0 {
		return ErrEmptyTitle
	}
	n.Title = title
	n.Description = description
	n.UpdatedAt = time.Now()
	return nil
}

func (n *Note) markAsCompleted() error {
	if n.Completed {
		return ErrNotIsAlreadyCompleted
	}
	n.Completed = true
	n.UpdatedAt = time.Now()
	return nil
}

func (n *Note) markAsInProgress() error {
	if !n.Completed {
		return ErrNotIsAlreadyInProgress
	}
	n.Completed = false
	n.UpdatedAt = time.Now()
	return nil
}

func (n *Note) copy() Note {
	return Note{
		ID:          uuid.Must(uuid.NewV4()),
		Title:       n.Title,
		Description: n.Description,
		UserID:      n.UserID,
		CreatedAt:   time.Now(),
	}
}
