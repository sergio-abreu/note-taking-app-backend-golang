package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type MarkNoteAsInProgressResponse struct {
	NoteID    uuid.UUID `json:"note_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a CommandApplication) MarkNoteAsInProgress(ctx context.Context, userID, noteID string) (MarkNoteAsInProgressResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return MarkNoteAsInProgressResponse{}, err
	}

	note, err := a.notesRepo.FindNote(ctx, userID, noteID)
	if err != nil {
		return MarkNoteAsInProgressResponse{}, err
	}

	err = user.MarkNoteAsInProgress(&note)
	if err != nil {
		return MarkNoteAsInProgressResponse{}, err
	}

	err = a.notesRepo.MarkAsInProgress(ctx, note)
	if err != nil {
		return MarkNoteAsInProgressResponse{}, err
	}

	return MarkNoteAsInProgressResponse{
		NoteID:    note.ID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}
