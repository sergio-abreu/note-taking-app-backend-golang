package notes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type ScheduleReminderRequest struct {
	StartDate  string `json:"start_date"`
	StartTime  string `json:"start_time"`
	Timezone   string `json:"timezone"`
	Interval   string `json:"interval"`
	WeekDays   string `json:"week_days"`
	EndsAfterN uint   `json:"ends_after_n"`
	EndsAt     string `json:"ends_at"`
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

	reminder, err := user.ScheduleAReminder(
		note,
		r.StartDate,
		r.StartTime,
		r.Timezone,
		r.Interval,
		r.WeekDays,
		r.EndsAt,
		r.EndsAfterN,
	)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	err = a.notesRepo.ScheduleReminder(ctx, reminder)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	err = a.cron.CreateCron(ctx, reminder)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	return ScheduleReminderResponse{
		ReminderID: reminder.ID,
		CreatedAt:  reminder.CreatedAt,
		UpdatedAt:  reminder.UpdatedAt,
	}, nil
}
