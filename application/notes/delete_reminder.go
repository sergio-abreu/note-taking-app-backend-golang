package notes

import (
	"context"
)

func (a CommandApplication) DeleteReminder(ctx context.Context, userID, noteID, reminderID string) error {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	reminder, err := a.notesRepo.FindReminder(ctx, userID, noteID, reminderID)
	if err != nil {
		return err
	}

	err = user.DeleteReminder(reminder)
	if err != nil {
		return err
	}

	err = a.notesRepo.DeleteReminder(ctx, reminder)
	if err != nil {
		return err
	}

	err = a.cron.DeleteCron(ctx, reminder)
	if err != nil {
		return err
	}

	return nil
}
