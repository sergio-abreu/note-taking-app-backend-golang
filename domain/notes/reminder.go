package notes

import (
	"errors"
	"time"

	"github.com/adhocore/gronx"
	"github.com/gofrs/uuid"
)

var (
	ErrInvalidEndsAt                   = errors.New("invalid ends at date format")
	ErrInvalidCronExpression           = errors.New("invalid cron expression")
	ErrExceededMinimumTimeInterval     = errors.New("the minimum time interval for reminders is 24 hours")
	ErrCannotConfigureEndsAtAndRepeats = errors.New("cannot configure endsAt and repeats at the same time")
)

func newReminder(noteID, userID uuid.UUID, cronExpression, rawEndsAt string, repeats uint) (Reminder, error) {
	var endsAt time.Time
	var err error
	if len(rawEndsAt) > 0 {
		endsAt, err = time.Parse(time.RFC3339, rawEndsAt)
		if err != nil {
			return Reminder{}, ErrInvalidEndsAt
		}
	}
	if !endsAt.IsZero() && repeats > 0 {
		return Reminder{}, ErrCannotConfigureEndsAtAndRepeats
	}
	g := gronx.New()
	if !g.IsValid(cronExpression) {
		return Reminder{}, ErrInvalidCronExpression
	}
	nextTime, _ := gronx.NextTick(cronExpression, false)
	nextTimeAfter, _ := gronx.NextTickAfter(cronExpression, nextTime, false)
	if nextTimeAfter.Sub(nextTime) < 24*time.Hour {
		return Reminder{}, ErrExceededMinimumTimeInterval
	}
	return Reminder{
		ID:             uuid.Must(uuid.NewV4()),
		NoteID:         noteID,
		UserID:         userID,
		CronExpression: cronExpression,
		EndsAt:         endsAt,
		Repeats:        repeats,
	}, nil
}

type Reminder struct {
	ID             uuid.UUID
	NoteID         uuid.UUID
	UserID         uuid.UUID
	CronExpression string
	EndsAt         time.Time
	Repeats        uint
}
