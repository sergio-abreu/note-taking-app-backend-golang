package notes

import (
	"context"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
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

	reminder, err := a.notesRepo.FindReminderByNoteID(ctx, userID, noteID)
	if err != nil && err != notes.ErrReminderNotFound {
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

	err = a.cron.DeleteCron(ctx, reminder)
	if err != nil {
		return err
	}

	return nil
}
