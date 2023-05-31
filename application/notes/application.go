package notes

import (
	"gorm.io/gorm"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewCommandApplication(notesRepo notes.Repository, cron notes.Cron) CommandApplication {
	return CommandApplication{notesRepo: notesRepo, cron: cron}
}

type CommandApplication struct {
	notesRepo notes.Repository
	cron      notes.Cron
}

func NewQueryApplication(db *gorm.DB) QueryApplication {
	return QueryApplication{db: db}
}

type QueryApplication struct {
	db *gorm.DB
}
