package notes

import (
	"context"
)

func (a CommandApplication) DeleteNote(ctx context.Context, userID, noteID string) error {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	note, err := a.notesRepo.FindNote(ctx, userID, noteID)
	if err != nil {
		return err
	}

	err = user.DeleteNote(note)
	if err != nil {
		return err
	}

	err = a.notesRepo.DeleteNote(ctx, note)
	if err != nil {
		return err
	}

	return nil
}
