package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type MarkNoteAsInProgressResponse struct {
	NoteID uuid.UUID `json:"note_id,omitempty"`
}

func (a Application) MarkNoteAsInProgress(ctx context.Context, userID, noteID string) (MarkNoteAsInProgressResponse, error) {
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
		NoteID: note.ID,
	}, nil
}
