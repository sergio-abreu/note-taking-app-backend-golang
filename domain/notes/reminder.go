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
	endsAt, err := validateEndsAt(cronExpression, rawEndsAt, repeats)
	if err != nil {
		return Reminder{}, err
	}
	err = validateCronExpression(cronExpression)
	if err != nil {
		return Reminder{}, err
	}
	return Reminder{
		ID:             uuid.Must(uuid.NewV4()),
		NoteID:         noteID,
		UserID:         userID,
		CronExpression: cronExpression,
		EndsAt:         endsAt,
		CreatedAt:      time.Now(),
	}, nil
}

type Reminder struct {
	ID             uuid.UUID
	NoteID         uuid.UUID
	UserID         uuid.UUID
	CronExpression string
	EndsAt         time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (r *Reminder) reschedule(cronExpression, rawEndsAt string, repeats uint) error {
	err := validateCronExpression(cronExpression)
	if err != nil {
		return err
	}
	endsAt, err := validateEndsAt(cronExpression, rawEndsAt, repeats)
	if err != nil {
		return err
	}
	r.EndsAt = endsAt
	r.CronExpression = cronExpression
	r.UpdatedAt = time.Now()
	return nil
}

func validateEndsAt(cronExpression, rawEndsAt string, repeats uint) (endsAt time.Time, err error) {
	if len(rawEndsAt) == 0 && repeats == 0 {
		return
	}
	if len(rawEndsAt) > 0 && repeats > 0 {
		return time.Time{}, ErrCannotConfigureEndsAtAndRepeats
	}
	if repeats > 0 {
		endsAt = time.Now()
		for i := uint(0); i < repeats; i++ {
			endsAt, _ = gronx.NextTickAfter(cronExpression, endsAt, false)
		}
		return
	}
	endsAt, err = time.Parse(time.RFC3339, rawEndsAt)
	if err != nil {
		return endsAt, ErrInvalidEndsAt
	}
	return
}

func validateCronExpression(cronExpression string) error {
	g := gronx.New()
	if !g.IsValid(cronExpression) {
		return ErrInvalidCronExpression
	}
	nextTime, _ := gronx.NextTick(cronExpression, false)
	nextTimeAfter, _ := gronx.NextTickAfter(cronExpression, nextTime, false)
	if nextTimeAfter.Sub(nextTime) < 24*time.Hour {
		return ErrExceededMinimumTimeInterval
	}
	return nil
}
