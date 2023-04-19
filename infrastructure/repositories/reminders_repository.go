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
		err = r.db.WithContext(ctx).
			Exec(sql, reminder.NoteID, reminder.CronExpression).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r RemindersRepository) RescheduleReminder(ctx context.Context, reminder notes.Reminder) error {
	//TODO implement me
	panic("implement me")
}

func (r RemindersRepository) DeleteReminder(ctx context.Context, reminder notes.Reminder) error {
	//TODO implement me
	panic("implement me")
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
