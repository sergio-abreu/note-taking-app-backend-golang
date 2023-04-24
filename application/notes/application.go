package notes

import (
	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewApplication(notesRepo notes.NotesRepository) Application {
	return Application{notesRepo: notesRepo}
}

type Application struct {
	notesRepo notes.NotesRepository
}
