package repositories

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewNotesRepository(db *gorm.DB) NotesRepository {
	return NotesRepository{db: db}
}

type NotesRepository struct {
	db *gorm.DB
}

func (n NotesRepository) FindUser(ctx context.Context, userID string) (notes.User, error) {
	id, err := parseUserID(userID)
	if err != nil {
		return notes.User{}, err
	}

	var user notes.User
	err = n.db.WithContext(ctx).
		Table("users").
		First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return notes.User{}, notes.ErrUserNotFound
	}
	if err != nil {
		return notes.User{}, err
	}

	return user, nil
}

func (n NotesRepository) CreateUser(ctx context.Context, user notes.User) error {
	err := n.db.WithContext(ctx).
		Table("users").
		Select("id", "name", "email", "created_at", "updated_at").
		Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (n NotesRepository) FindNote(ctx context.Context, rawUserID, rawNoteID string) (notes.Note, error) {
	userID, err := parseUserID(rawUserID)
	if err != nil {
		return notes.Note{}, err
	}
	noteID, err := parseNoteID(rawNoteID)
	if err != nil {
		return notes.Note{}, err
	}

	var note notes.Note
	err = n.db.WithContext(ctx).
		Table("notes").
		First(&note, "user_id = ? AND id = ?", userID, noteID).Error
	if err == gorm.ErrRecordNotFound {
		return notes.Note{}, notes.ErrNoteNotFound
	}
	if err != nil {
		return notes.Note{}, err
	}

	return note, nil
}

func (n NotesRepository) CreateNote(ctx context.Context, note notes.Note) error {
	err := n.db.WithContext(ctx).
		Table("notes").
		Select("id", "title", "description", "user_id", "created_at", "updated_at").
		Create(&note).Error
	if err != nil {
		return err
	}

	return nil
}

func (n NotesRepository) EditNote(ctx context.Context, note notes.Note) error {
	err := n.db.WithContext(ctx).
		Table("notes").
		Select("title", "description", "updated_at").
		Where("id = ?", note.ID).
		Updates(&note).Error
	if err == gorm.ErrRecordNotFound {
		return notes.ErrNoteNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (n NotesRepository) DeleteNote(ctx context.Context, note notes.Note) error {
	err := n.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var reminder notes.Reminder
		err := n.db.WithContext(ctx).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Table("reminders").
			Where("note_id = ?", note.ID).
			Delete(&reminder).Error
		if err != nil {
			return err
		}

		err = n.db.WithContext(ctx).
			Table("notes").
			Delete(&note).Error
		if err == gorm.ErrRecordNotFound {
			return notes.ErrNoteNotFound
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (n NotesRepository) MarkAsComplete(ctx context.Context, note notes.Note) error {
	err := n.db.WithContext(ctx).
		Table("notes").
		Select("completed", "updated_at").
		Where("id = ?", note.ID).
		Updates(&note).Error
	if err == gorm.ErrRecordNotFound {
		return notes.ErrNoteNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (n NotesRepository) MarkAsInProgress(ctx context.Context, note notes.Note) error {
	err := n.db.WithContext(ctx).
		Table("notes").
		Select("completed", "updated_at").
		Where("id = ?", note.ID).
		Updates(&note).Error
	if err == gorm.ErrRecordNotFound {
		return notes.ErrNoteNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (n NotesRepository) FindReminder(ctx context.Context, rawUserID, rawNoteID, rawReminderID string) (notes.Reminder, error) {
	userID, err := parseUserID(rawUserID)
	if err != nil {
		return notes.Reminder{}, err
	}
	noteID, err := parseNoteID(rawNoteID)
	if err != nil {
		return notes.Reminder{}, err
	}
	reminderID, err := parseReminderID(rawReminderID)
	if err != nil {
		return notes.Reminder{}, err
	}

	var reminder notes.Reminder
	err = n.db.WithContext(ctx).
		Table("reminders").
		First(&reminder, "user_id = ? AND note_id = ? AND id = ?", userID, noteID, reminderID).Error
	if err == gorm.ErrRecordNotFound {
		return notes.Reminder{}, notes.ErrReminderNotFound
	}
	if err != nil {
		return notes.Reminder{}, err
	}

	return reminder, nil
}

func (n NotesRepository) FindReminderByNoteID(ctx context.Context, rawUserID, rawNoteID string) (notes.Reminder, error) {
	userID, err := parseUserID(rawUserID)
	if err != nil {
		return notes.Reminder{}, err
	}
	noteID, err := parseNoteID(rawNoteID)
	if err != nil {
		return notes.Reminder{}, err
	}

	var reminder notes.Reminder
	err = n.db.WithContext(ctx).
		Table("reminders").
		First(&reminder, "user_id = ? AND note_id = ?", userID, noteID).Error
	if err == gorm.ErrRecordNotFound {
		return notes.Reminder{}, notes.ErrReminderNotFound
	}
	if err != nil {
		return notes.Reminder{}, err
	}

	return reminder, nil
}

func (n NotesRepository) ScheduleReminder(ctx context.Context, reminder notes.Reminder) error {
	return n.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		fields := []string{"id", "note_id", "user_id", "start_date", "start_time", "timezone", "interval", "ends_after_n", "created_at", "updated_at"}
		if !reminder.EndsAt.IsZero() {
			fields = append(fields, "ends_at")
		}
		if reminder.Interval == notes.Weekly && len(reminder.WeekDays) != 0 {
			fields = append(fields, "week_days")
		}
		err := tx.WithContext(ctx).
			Table("reminders").
			Select(fields).
			Create(&reminder).Error
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.ConstraintName == "unique_notes_idx" {
			return notes.ErrOnlyOneReminderAllowed
		}
		if err != nil {
			return err
		}

		return nil
	})
}

func (n NotesRepository) RescheduleReminder(ctx context.Context, reminder notes.Reminder) error {
	return n.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		fields := []string{"start_date", "start_time", "timezone", "interval", "ends_after_n", "updated_at"}
		if !reminder.EndsAt.IsZero() {
			fields = append(fields, "ends_at")
		}
		if reminder.Interval == notes.Weekly && len(reminder.WeekDays) != 0 {
			fields = append(fields, "week_days")
		}
		return tx.WithContext(ctx).
			Table("reminders").
			Select(fields).
			Where("id = ?", reminder.ID).
			Updates(&reminder).Error
	})
}

func (n NotesRepository) DeleteReminder(ctx context.Context, reminder notes.Reminder) error {
	return n.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).
			Table("reminders").
			Where("id = ?", reminder.ID).
			Delete(&reminder).Error
	})
}

func parseUserID(userID string) (uuid.UUID, error) {
	if len(userID) == 0 {
		return uuid.Nil, notes.ErrEmptyUserID
	}
	id, err := uuid.FromString(userID)
	if err != nil {
		return uuid.Nil, notes.ErrInvalidUserIDFormat
	}
	return id, nil
}

func parseNoteID(noteID string) (uuid.UUID, error) {
	if len(noteID) == 0 {
		return uuid.Nil, notes.ErrEmptyNoteID
	}
	id, err := uuid.FromString(noteID)
	if err != nil {
		return uuid.Nil, notes.ErrInvalidNoteIDFormat
	}
	return id, nil
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
