package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type MarkNoteAsCompleteResponse struct {
	NoteID uuid.UUID `json:"note_id,omitempty"`
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
		NoteID: note.ID,
	}, nil
}
