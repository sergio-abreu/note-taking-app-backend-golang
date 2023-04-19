package repositories

import (
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewNotesRepository(db *gorm.DB) NotesRepository {
	return NotesRepository{db: db}
}

type NotesRepository struct {
	db *gorm.DB
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
		Select("id", "title", "description", "user_id", "created_at").
		Omit("updated_at").
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
	//TODO implement me
	panic("implement me")
}

func (n NotesRepository) MarkAsComplete(ctx context.Context, note notes.Note) error {
	//TODO implement me
	panic("implement me")
}

func (n NotesRepository) MarkAsInProgress(ctx context.Context, note notes.Note) error {
	//TODO implement me
	panic("implement me")
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
