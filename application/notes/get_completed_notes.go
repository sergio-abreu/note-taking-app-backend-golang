package notes

import (
	"context"
)

func (q QueryApplication) GetCompletedNotes(ctx context.Context, userID string) ([]Note, error) {
	var noteList []Note
	err := q.db.WithContext(ctx).
		Table("notes n").
		Preload("Reminder").
		Where("n.user_id = ? AND n.completed = true", userID).
		Find(&noteList).Error
	return noteList, err
}
