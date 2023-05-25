package notes

import (
	"context"
	"fmt"
	"time"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

type Note struct {
	notes.Note
	Reminder *Reminder `json:"reminder"`
}

type Reminder struct {
	notes.Reminder
	StartDate Date `json:"start_date"`
	EndsAt    Date `json:"ends_at"`
}

type Date time.Time

func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	date := time.Time(*d)
	if date.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%q", date.Format(time.DateOnly))), nil
}

func (q QueryApplication) GetNotes(ctx context.Context, userID string) ([]Note, error) {
	var noteList []Note
	err := q.db.WithContext(ctx).
		Table("notes n").
		Preload("Reminder").
		Where("n.user_id = ? AND n.completed = false", userID).
		Order("created_at DESC").
		Find(&noteList).Error
	return noteList, err
}
