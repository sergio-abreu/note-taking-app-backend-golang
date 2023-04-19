package repositories

import (
	"context"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return UsersRepository{db: db}
}

type UsersRepository struct {
	db *gorm.DB
}

func (u UsersRepository) FindUser(ctx context.Context, userID string) (notes.User, error) {
	id, err := parseUserID(userID)
	if err != nil {
		return notes.User{}, err
	}

	var user notes.User
	err = u.db.WithContext(ctx).
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

func (u UsersRepository) CreateUser(ctx context.Context, user notes.User) error {
	err := u.db.WithContext(ctx).
		Table("users").
		Select("id", "name", "email", "created_at").
		Omit("updated_at").
		Create(&user).Error
	if err != nil {
		return err
	}

	return nil
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
