package notes

import (
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewCommandApplication(notesRepo notes.Repository) CommandApplication {
	return CommandApplication{notesRepo: notesRepo}
}

type CommandApplication struct {
	notesRepo notes.Repository
}
