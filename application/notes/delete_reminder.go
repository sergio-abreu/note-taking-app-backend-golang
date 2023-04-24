package notes

import (
	"context"
)

func (a Application) DeleteReminder(ctx context.Context, userID, reminderID string) error {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return err
	}

	reminder, err := a.notesRepo.FindReminder(ctx, userID, reminderID)
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

	return nil
}
