package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type RescheduleReminderRequest struct {
	CronExpression string `json:"cron_expression,omitempty"`
	EndsAt         string `json:"ends_at,omitempty"`
	Repeats        uint   `json:"repeats,omitempty"`
}

type RescheduleReminderResponse struct {
	ReminderID uuid.UUID `json:"reminder_id,omitempty"`
}

func (a Application) RescheduleReminder(ctx context.Context, userID, noteID, reminderID string, r RescheduleReminderRequest) (RescheduleReminderResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	reminder, err := a.notesRepo.FindReminder(ctx, userID, noteID, reminderID)
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
