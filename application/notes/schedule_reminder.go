package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type ScheduleReminderRequest struct {
	CronExpression string `json:"cron_expression"`
	EndsAt         string `json:"ends_at"`
	Repeats        uint   `json:"repeats"`
}

type ScheduleReminderResponse struct {
	ReminderID uuid.UUID `json:"reminder_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (a CommandApplication) ScheduleReminder(ctx context.Context, userID, noteID string, r ScheduleReminderRequest) (ScheduleReminderResponse, error) {
	user, err := a.notesRepo.FindUser(ctx, userID)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	note, err := a.notesRepo.FindNote(ctx, userID, noteID)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	reminder, err := user.ScheduleAReminder(note, r.CronExpression, r.EndsAt, r.Repeats)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	err = a.notesRepo.ScheduleReminder(ctx, reminder)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	return ScheduleReminderResponse{
		ReminderID: reminder.ID,
		CreatedAt:  reminder.CreatedAt,
		UpdatedAt:  reminder.UpdatedAt,
	}, nil
}
