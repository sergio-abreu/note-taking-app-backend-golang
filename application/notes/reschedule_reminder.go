package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type RescheduleReminderRequest struct {
	CronExpression string `json:"cron_expression"`
	EndsAt         string `json:"ends_at"`
	Repeats        uint   `json:"repeats"`
}

type RescheduleReminderResponse struct {
	ReminderID uuid.UUID `json:"reminder_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (a CommandApplication) RescheduleReminder(ctx context.Context, userID, noteID, reminderID string, r RescheduleReminderRequest) (RescheduleReminderResponse, error) {
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
		CreatedAt:  reminder.CreatedAt,
		UpdatedAt:  reminder.UpdatedAt,
	}, nil
}
