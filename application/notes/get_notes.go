package notes

import (
	"context"
)

func (q QueryApplication) GetNotes(ctx context.Context, userID string) ([]Note, error) {
	var noteList []Note
	err := q.db.WithContext(ctx).
		Table("notes n").
		Preload("Reminder").
		Where("n.user_id = ?", userID).
		Order("created_at DESC").
		Find(&noteList).Error
	return noteList, err
}
