package notes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adhocore/gronx"
	"github.com/gofrs/uuid"
)

var (
	ErrInvalidStartDate                   = errors.New("invalid start date format")
	ErrInvalidStartTime                   = errors.New("invalid start time format")
	ErrInvalidTimezone                    = errors.New("invalid timezone format")
	ErrInvalidInterval                    = errors.New("invalid interval")
	ErrInvalidWeekDays                    = errors.New("invalid week days")
	ErrInvalidEndsAt                      = errors.New("invalid ends at date format")
	ErrCannotConfigureMultipleTermination = errors.New("cannot configure multiple termination")
)

type Interval string

const (
	Daily   Interval = "Daily"
	Weekly  Interval = "Weekly"
	Monthly Interval = "Monthly"
	Yearly  Interval = "Yearly"
)

func newReminder(
	noteID uuid.UUID,
	userID uuid.UUID,
	rawStartDate string,
	startTime string,
	timezone string,
	rawInterval string,
	rawWeekDays string,
	rawEndsAt string,
	endsAfterN uint,
) (Reminder, error) {
	startDate, err := parseStartDate(rawStartDate)
	if err != nil {
		return Reminder{}, err
	}
	err = validateStartTime(startTime)
	if err != nil {
		return Reminder{}, err
	}
	err = validateTimezone(timezone)
	if err != nil {
		return Reminder{}, err
	}
	interval, err := parseInterval(rawInterval)
	if err != nil {
		return Reminder{}, err
	}
	weekDays, err := parseWeekDays(interval, rawWeekDays)
	if err != nil {
		return Reminder{}, err
	}
	endsAt, err := parseAndValidateEndsAt(rawEndsAt, endsAfterN)
	if err != nil {
		return Reminder{}, err
	}
	now := time.Now()
	return Reminder{
		ID:         uuid.Must(uuid.NewV4()),
		NoteID:     noteID,
		UserID:     userID,
		StartDate:  startDate,
		StartTime:  startTime,
		Timezone:   timezone,
		Interval:   interval,
		WeekDays:   weekDays,
		EndsAfterN: endsAfterN,
		EndsAt:     endsAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

type Reminder struct {
	ID         uuid.UUID `json:"id"`
	NoteID     uuid.UUID `json:"note_id"`
	UserID     uuid.UUID `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	StartTime  string    `json:"start_time"`
	Timezone   string    `json:"timezone"`
	Interval   Interval  `json:"every"`
	WeekDays   string    `json:"week_days"`
	EndsAfterN uint      `json:"ends_after_n"`
	EndsAt     time.Time `json:"ends_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *Reminder) reschedule(
	rawStartDate string,
	startTime string,
	timezone string,
	rawInterval string,
	rawWeekDays string,
	rawEndsAt string,
	endsAfterN uint,
) error {
	startDate, err := parseStartDate(rawStartDate)
	if err != nil {
		return err
	}
	err = validateStartTime(startTime)
	if err != nil {
		return err
	}
	err = validateTimezone(timezone)
	if err != nil {
		return err
	}
	interval, err := parseInterval(rawInterval)
	if err != nil {
		return err
	}
	weekDays, err := parseWeekDays(interval, rawWeekDays)
	if err != nil {
		return err
	}
	endsAt, err := parseAndValidateEndsAt(rawEndsAt, endsAfterN)
	if err != nil {
		return err
	}
	r.StartDate = startDate
	r.StartTime = startTime
	r.Timezone = timezone
	r.Interval = interval
	r.WeekDays = weekDays
	r.EndsAt = endsAt
	r.EndsAfterN = endsAfterN
	r.UpdatedAt = time.Now()
	return nil
}

func (r *Reminder) ParseCron() string {
	loc, _ := time.LoadLocation(r.Timezone)
	refDate := r.StartDate
	refTime := strings.Split(r.StartTime, ":")
	hour, _ := strconv.Atoi(refTime[0])
	minutes, _ := strconv.Atoi(refTime[1])
	startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), hour, minutes, 0, 0, loc).UTC()
	day := "*"
	month := "*"
	week := "*"
	switch r.Interval {
	case Daily:
	case Weekly:
		if len(r.WeekDays) != 0 {
			week = r.WeekDays
		}
	case Monthly:
		day = strconv.Itoa(startDate.Day())
	case Yearly:
		day = strconv.Itoa(startDate.Day())
		month = strconv.Itoa(int(startDate.Month()))
	}
	return fmt.Sprintf("%d %d %s %s %s", minutes, hour, day, month, week)
}

func (r *Reminder) ParseEndsAt(cronExpression string) string {
	if r.EndsAfterN > 0 {
		endsAt := time.Now().UTC()
		for i := uint(0); i < r.EndsAfterN; i++ {
			endsAt, _ = gronx.NextTickAfter(cronExpression, endsAt, false)
		}
		return endsAt.Format(time.RFC3339)
	}
	if !r.EndsAt.IsZero() {
		return r.EndsAt.Format(time.DateOnly)
	}
	return ""
}

func parseStartDate(rawStartDate string) (time.Time, error) {
	startDate, err := time.Parse(time.DateOnly, rawStartDate)
	if err != nil {
		return time.Time{}, ErrInvalidStartDate
	}
	return startDate, nil
}

func validateStartTime(rawStartTime string) error {
	if _, err := time.Parse("15:04", rawStartTime); err != nil {
		return ErrInvalidStartTime
	}
	return nil
}

func validateTimezone(timezone string) error {
	if _, err := time.LoadLocation(timezone); err != nil {
		return ErrInvalidTimezone
	}
	return nil
}

func parseInterval(rawInterval string) (Interval, error) {
	switch Interval(rawInterval) {
	case Daily:
		return Daily, nil
	case Weekly:
		return Weekly, nil
	case Monthly:
		return Monthly, nil
	case Yearly:
		return Yearly, nil
	}
	return "", ErrInvalidInterval
}

func parseWeekDays(interval Interval, weekDays string) (string, error) {
	if interval != Weekly || len(weekDays) == 0 {
		return "", nil
	}
	days := strings.Split(weekDays, ",")
	repeat := map[string]bool{}
	for _, day := range days {
		intDay, err := strconv.Atoi(day)
		if err != nil || repeat[day] || intDay < 1 || intDay > 7 {
			return "", ErrInvalidWeekDays
		}
		repeat[day] = true
	}
	return weekDays, nil
}

func parseAndValidateEndsAt(rawEndsAt string, repeats uint) (endsAt time.Time, err error) {
	if len(rawEndsAt) == 0 && repeats == 0 {
		return
	}
	if len(rawEndsAt) > 0 && repeats > 0 {
		return time.Time{}, ErrCannotConfigureMultipleTermination
	}
	if len(rawEndsAt) > 0 {
		endsAt, err = time.Parse(time.DateOnly, rawEndsAt)
		if err != nil {
			return endsAt, ErrInvalidEndsAt
		}
	}
	return
}
