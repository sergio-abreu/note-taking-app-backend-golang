package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/sergio-abreu/note-taking-app-backend-golang/application/notes"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/cron"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/repositories"
)

type WebServer struct {
	command notes.CommandApplication
	query   notes.QueryApplication
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err.Error())
	}
}

func run() error {
	db, err := repositories.NewGormDBFromEnv()
	if err != nil {
		return err
	}

	db = db.Debug()
	notesRepo := repositories.NewNotesRepository(db)
	localCron := cron.NewLocalCron("/tmp")
	command := notes.NewCommandApplication(notesRepo, localCron)
	query := notes.NewQueryApplication(db)
	server := WebServer{command: command, query: query}

	r := gin.Default()
	r.Use(cors.Default())
	g := r.Group("/api/v1/:userID/notes")
	g.GET("", server.GetNotes)
	g.GET("/in-progress", server.GetInProgressNotes)
	g.GET("/completed", server.GetCompletedNotes)
	g.POST("", server.CreateNote)
	g.PATCH("/:noteID", server.EditNote)
	g.DELETE("/:noteID", server.DeleteNote)
	g.POST("/:noteID/copy", server.CopyNote)
	g.PUT("/:noteID/complete", server.MarkNoteAsComplete)
	g.PUT("/:noteID/in-progress", server.MarkNoteAsInProgress)
	g.POST("/:noteID/reminders", server.ScheduleReminder)
	g.PATCH("/:noteID/reminders/:reminderID", server.RescheduleReminder)
	g.DELETE("/:noteID/reminders/:reminderID", server.DeleteReminder)

	if err := r.Run(":8080"); err != nil {
		return err
	}

	return nil
}
