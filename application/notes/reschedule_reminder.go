package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type RescheduleReminderRequest struct {
	StartDate  string `json:"start_date"`
	StartTime  string `json:"start_time"`
	Timezone   string `json:"timezone"`
	Interval   string `json:"interval"`
	WeekDays   string `json:"week_days"`
	EndsAfterN uint   `json:"ends_after_n"`
	EndsAt     string `json:"ends_at"`
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

	err = user.RescheduleAReminder(&reminder, r.StartDate, r.StartTime, r.Timezone, r.Interval, r.WeekDays, r.EndsAt, r.EndsAfterN)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	err = a.notesRepo.RescheduleReminder(ctx, reminder)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	err = a.cron.CreateCron(ctx, reminder)
	if err != nil {
		return RescheduleReminderResponse{}, err
	}

	return RescheduleReminderResponse{
		ReminderID: reminder.ID,
		CreatedAt:  reminder.CreatedAt,
		UpdatedAt:  reminder.UpdatedAt,
	}, nil
}
