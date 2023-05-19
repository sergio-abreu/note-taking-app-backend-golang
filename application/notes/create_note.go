package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type CreateNoteRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateNoteResponse struct {
	NoteID    uuid.UUID `json:"note_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
		NoteID:    note.ID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}
