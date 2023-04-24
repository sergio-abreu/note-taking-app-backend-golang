package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type EditNoteRequest struct {
	Title       string
	Description string
}

type EditNoteResponse struct {
	NoteID uuid.UUID
}

func (a Application) EditNote(ctx context.Context, userID, noteID string, r EditNoteRequest) (EditNoteResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return EditNoteResponse{}, err
	}

	note, err := a.notesRepo.FindNote(ctx, userID, noteID)
	if err != nil {
		return EditNoteResponse{}, err
	}

	err = user.EditNote(&note, r.Title, r.Description)
	if err != nil {
		return EditNoteResponse{}, err
	}

	err = a.notesRepo.EditNote(ctx, note)
	if err != nil {
		return EditNoteResponse{}, err
	}

	return EditNoteResponse{
		NoteID: note.ID,
	}, nil
}
