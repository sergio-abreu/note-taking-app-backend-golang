package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type CreateNoteRequest struct {
	Title       string
	Description string
}

type CreateNoteResponse struct {
	NoteID uuid.UUID
}

func (a Application) CreateNote(ctx context.Context, userID string, r CreateNoteRequest) (CreateNoteResponse, error) {
	user, err := a.usersRepo.FindUser(ctx, userID)
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
