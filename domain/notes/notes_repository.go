package notes

import (
	"context"
	"errors"
)

var (
	ErrEmptyNoteID         = errors.New("empty note id")
	ErrInvalidNoteIDFormat = errors.New("invalid note id format")
	ErrNoteNotFound        = errors.New("note not found")
)

type NotesRepository interface {
	FindNote(ctx context.Context, userID, noteID string) (Note, error)
	CreateNote(ctx context.Context, note Note) error
	EditNote(ctx context.Context, note Note) error
	DeleteNote(ctx context.Context, note Note) error
	MarkAsComplete(ctx context.Context, note Note) error
	MarkAsInProgress(ctx context.Context, note Note) error
}
