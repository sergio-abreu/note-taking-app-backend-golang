package emailer

import (
	"context"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/emailer"
	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewApplication(usersRepo notes.UsersRepository, notesRepo notes.NotesRepository, emailer emailer.Emailer) Application {
	return Application{usersRepo: usersRepo, notesRepo: notesRepo, emailer: emailer}
}

type Application struct {
	usersRepo notes.UsersRepository
	notesRepo notes.NotesRepository
	emailer   emailer.Emailer
}

type SendReminderEmailRequest struct {
	ReminderID string `json:"reminder_id,omitempty"`
	NoteID     string `json:"note_id,omitempty"`
	UserID     string `json:"user_id,omitempty"`
}

func (a Application) SendReminderEmail(ctx context.Context, r SendReminderEmailRequest) error {
	user, err := a.usersRepo.FindUser(ctx, r.UserID)
	if err != nil {
		return err
	}

	note, err := a.notesRepo.FindNote(ctx, r.UserID, r.NoteID)
	if err != nil {
		return err
	}

	err = a.emailer.SendNoteReminder(ctx, user.Email, note)
	if err != nil {
		return err
	}

	return nil
}
