package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type EditNoteRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type EditNoteResponse struct {
	NoteID    uuid.UUID `json:"note_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a CommandApplication) EditNote(ctx context.Context, userID, noteID string, r EditNoteRequest) (EditNoteResponse, error) {
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
		NoteID:    note.ID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}
