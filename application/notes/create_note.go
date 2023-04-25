package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type CreateNoteRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateNoteResponse struct {
	NoteID uuid.UUID `json:"note_id,omitempty"`
}

func (a CommandApplication) CreateNote(ctx context.Context, userID string, r CreateNoteRequest) (CreateNoteResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return CreateNoteResponse{}, err
	}

	note, err := user.CreateNote(r.Title, r.Description)
	if err != nil {
		return CreateNoteResponse{}, err
	}

	err = a.notesRepo.CreateNote(ctx, note)
	if err != nil {
		return CreateNoteResponse{}, err
	}

	return CreateNoteResponse{
		NoteID: note.ID,
	}, nil
}
