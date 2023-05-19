package notes

import (
	"context"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

type Note struct {
	notes.Note
	Reminder *notes.Reminder `json:"reminder"`
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
