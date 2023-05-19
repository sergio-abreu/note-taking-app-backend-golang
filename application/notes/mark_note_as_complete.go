package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type MarkNoteAsCompleteResponse struct {
	NoteID    uuid.UUID `json:"note_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a CommandApplication) MarkNoteAsComplete(ctx context.Context, userID, noteID string) (MarkNoteAsCompleteResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return MarkNoteAsCompleteResponse{}, err
	}

	note, err := a.notesRepo.FindNote(ctx, userID, noteID)
	if err != nil {
		return MarkNoteAsCompleteResponse{}, err
	}

	err = user.MarkNoteAsCompleted(&note)
	if err != nil {
		return MarkNoteAsCompleteResponse{}, err
	}

	err = a.notesRepo.MarkAsComplete(ctx, note)
	if err != nil {
		return MarkNoteAsCompleteResponse{}, err
	}

	return MarkNoteAsCompleteResponse{
		NoteID:    note.ID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}
