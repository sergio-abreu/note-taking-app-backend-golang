package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type CopyNoteRequest struct {
	Title       string
	Description string
}

type CopyNoteResponse struct {
	NoteID uuid.UUID
}

func (a Application) CopyNote(ctx context.Context, userID, noteID string) (CopyNoteResponse, error) {
	user, err := a.usersRepo.FindUser(ctx, userID)
	if err != nil {
		return CopyNoteResponse{}, err
	}

	note, err := a.notesRepo.FindNote(ctx, userID, noteID)
	if err != nil {
		return CopyNoteResponse{}, err
	}

	newNote, err := user.CopyNote(note)
	if err != nil {
		return CopyNoteResponse{}, err
	}

	err = a.notesRepo.CreateNote(ctx, newNote)
	if err != nil {
		return CopyNoteResponse{}, err
	}

	return CopyNoteResponse{
		NoteID: newNote.ID,
	}, nil
}
