package emailer

import (
	"context"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

//go:generate mockgen -source emailer.go -destination emailer_mock.go -package emailer
type Emailer interface {
	SendNoteReminder(ctx context.Context, to string, note notes.Note) error
}
