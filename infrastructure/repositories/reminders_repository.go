package repositories

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewRemindersRepository(db *gorm.DB) RemindersRepository {
	return RemindersRepository{db: db}
}

type RemindersRepository struct {
	db *gorm.DB
}

func (r RemindersRepository) FindReminder(ctx context.Context, rawUserID, rawReminderID string) (notes.Reminder, error) {
	userID, err := parseUserID(rawUserID)
	if err != nil {
		return notes.Reminder{}, err
	}
	reminderID, err := parseReminderID(rawReminderID)
	if err != nil {
		return notes.Reminder{}, err
	}

	var reminder notes.Reminder
	err = r.db.WithContext(ctx).
		Table("reminders").
		First(&reminder, "user_id = ? AND id = ?", userID, reminderID).Error
	if err == gorm.ErrRecordNotFound {
		return notes.Reminder{}, notes.ErrReminderNotFound
	}
	if err != nil {
		return notes.Reminder{}, err
	}

	return reminder, nil
}

func (r RemindersRepository) ScheduleReminder(ctx context.Context, reminder notes.Reminder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).
			Table("reminders").
			Select("id", "note_id", "user_id", "cron_expression", "ends_at", "repeats", "created_at").
			Omit("updated_at").
			Create(&reminder).Error
		if err != nil {
			return err
		}

		return r.scheduleCron(ctx, reminder)
	})
}

func (r RemindersRepository) RescheduleReminder(ctx context.Context, reminder notes.Reminder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).
			Table("reminders").
			Select("cron_expression", "ends_at", "repeats", "updated_at").
			Where("id = ?", reminder.ID).
			Updates(&reminder).Error
		if err != nil {
			return err
		}

		return r.scheduleCron(ctx, reminder)
	})
}

func (r RemindersRepository) DeleteReminder(ctx context.Context, reminder notes.Reminder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).
			Table("reminders").
			Where("id = ?", reminder.ID).
			Delete(&reminder).Error
		if err != nil {
			return err
		}

		return r.unscheduleCron(ctx, reminder.NoteID)
	})
}

func (r RemindersRepository) scheduleCron(ctx context.Context, reminder notes.Reminder) error {
	sql := fmt.Sprintf(`
			SELECT cron.schedule(?, ?,
				$$
				SELECT status
				from
				  http_post(
					'http://note-taking-app.com/v1/webhooks/reminders/%s',
					'{"reminder_id": "%s", "note_id": "%s", "user_id": "%s"}',
					'application/json'
				  )
				$$
			  );
		`, reminder.ID, reminder.ID, reminder.NoteID, reminder.UserID)
	return r.db.WithContext(ctx).Exec(sql, reminder.NoteID, reminder.CronExpression).Error
}

func (r RemindersRepository) unscheduleCron(ctx context.Context, noteID uuid.UUID) error {
	sql := "SELECT cron.unschedule(?);"
	return r.db.WithContext(ctx).Exec(sql, noteID).Error
}

func parseReminderID(reminderID string) (uuid.UUID, error) {
	if len(reminderID) == 0 {
		return uuid.Nil, notes.ErrEmptyReminderID
	}
	id, err := uuid.FromString(reminderID)
	if err != nil {
		return uuid.Nil, notes.ErrInvalidReminderIDFormat
	}
	return id, nil
}
