package notes

import (
	"context"

	"github.com/gofrs/uuid"
)

type ScheduleReminderRequest struct {
	CronExpression string
	EndsAt         string
	Repeats        uint
}

type ScheduleReminderResponse struct {
	ReminderID uuid.UUID
}

func (a Application) ScheduleReminder(ctx context.Context, userID, noteID string, r ScheduleReminderRequest) (ScheduleReminderResponse, error) {
	user, err := a.usersRepo.FindUser(ctx, userID)
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

	err = a.remindersRepo.ScheduleReminder(ctx, reminder)
	if err != nil {
		return ScheduleReminderResponse{}, err
	}

	return ScheduleReminderResponse{
		ReminderID: reminder.ID,
	}, nil
}
