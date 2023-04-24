package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type RescheduleReminderRequest struct {
	CronExpression string
	EndsAt         string
	Repeats        uint
}

type RescheduleReminderResponse struct {
	ReminderID uuid.UUID
}

func (a Application) RescheduleReminder(ctx context.Context, userID, reminderID string, r RescheduleReminderRequest) (RescheduleReminderResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	reminder, err := a.notesRepo.FindReminder(ctx, userID, reminderID)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	err = user.RescheduleAReminder(&reminder, r.CronExpression, r.EndsAt, r.Repeats)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	err = a.notesRepo.RescheduleReminder(ctx, reminder)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	return RescheduleReminderResponse{
		ReminderID: reminder.ID,
	}, nil
}
