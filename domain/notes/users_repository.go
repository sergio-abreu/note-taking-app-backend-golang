package notes

import (
	"context"
	"errors"
)

var (
	ErrEmptyUserID         = errors.New("empty user id")
	ErrInvalidUserIDFormat = errors.New("invalid user id format")
	ErrUserNotFound        = errors.New("user not found")
)

type UsersRepository interface {
	FindUser(ctx context.Context, userID string) (User, error)
	CreateUser(ctx context.Context, user User) error
}
