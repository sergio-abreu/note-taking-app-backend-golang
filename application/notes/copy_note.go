package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type CopyNoteRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type CopyNoteResponse struct {
	NoteID uuid.UUID `json:"note_id,omitempty"`
}

func (a Application) CopyNote(ctx context.Context, userID, noteID string) (CopyNoteResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
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
