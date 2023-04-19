package notes

import (
	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewApplication(usersRepo notes.UsersRepository, notesRepo notes.NotesRepository, remindersRepo notes.RemindersRepository) Application {
	return Application{usersRepo: usersRepo, notesRepo: notesRepo, remindersRepo: remindersRepo}
}

type Application struct {
	usersRepo     notes.UsersRepository
	notesRepo     notes.NotesRepository
	remindersRepo notes.RemindersRepository
}
