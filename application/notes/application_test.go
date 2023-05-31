package notes

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/cron"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/repositories"
)

func TestApplication(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()
	notesRepo, app, err := initializeApplication(t)
	g.Expect(err).Should(
		Not(HaveOccurred()))

	t.Run("Create a note successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"

		r, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"NoteID":    Not(Equal(uuid.Nil)),
			"CreatedAt": BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt": BeTemporally("~", time.Now(), time.Second),
		}))
		noteFromDB, err := notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(noteFromDB).Should(
			notes.BeANote(t, title, description, false, fakeUser.ID, time.Now(), time.Now()))
	})

	t.Run("Edit a note successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       "test title",
			Description: "test description",
		})
		g.Expect(err).Should(
			Not(HaveOccurred()))
		editedTitle := "edited test title"
		editedDescription := "edited test description"

		r, err := app.EditNote(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), EditNoteRequest{
			Title:       editedTitle,
			Description: editedDescription,
		})

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"NoteID":    Equal(r.NoteID),
			"CreatedAt": BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt": BeTemporally("~", time.Now(), time.Second),
		}))
		noteFromDB, err := notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(noteFromDB).Should(
			notes.BeANote(t, editedTitle, editedDescription, false, fakeUser.ID, time.Now(), time.Now()))
	})

	t.Run("Mark note as completed successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})
		g.Expect(err).Should(
			Not(HaveOccurred()))

		r, err := app.MarkNoteAsComplete(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String())

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"NoteID":    Equal(r.NoteID),
			"CreatedAt": BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt": BeTemporally("~", time.Now(), time.Second),
		}))
		noteFromDB, err := notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(noteFromDB).Should(
			notes.BeANote(t, title, description, true, fakeUser.ID, time.Now(), time.Now()))
	})

	t.Run("Mark note as in progress successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})
		g.Expect(err).Should(
			Not(HaveOccurred()))
		_, err = app.MarkNoteAsComplete(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))

		r, err := app.MarkNoteAsInProgress(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String())

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"NoteID":    Equal(r.NoteID),
			"CreatedAt": BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt": BeTemporally("~", time.Now(), time.Second),
		}))
		noteFromDB, err := notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(noteFromDB).Should(
			notes.BeANote(t, title, description, false, fakeUser.ID, time.Now(), time.Now()))
	})

	t.Run("Copy a note successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})

		r, err := app.CopyNote(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String())

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"NoteID":    Not(Equal(uuid.Nil)),
			"CreatedAt": BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt": BeTemporally("~", time.Now(), time.Second),
		}))
		noteFromDB, err := notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(noteFromDB).Should(
			notes.BeANote(t, title, description, false, fakeUser.ID, time.Now(), time.Now()))
	})

	t.Run("Delete a note successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		r, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})

		err = app.DeleteNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			MatchError(notes.ErrNoteNotFound))
	})

	t.Run("Schedule a reminder successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "20:33"
		interval := notes.Daily

		r, err := app.ScheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), ScheduleReminderRequest{
			StartDate: startDate.Format(time.DateOnly),
			StartTime: startTime,
			Timezone:  timezone,
			Interval:  string(interval),
		})

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"ReminderID": Not(Equal(uuid.Nil)),
			"CreatedAt":  BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt":  BeTemporally("~", time.Now(), time.Second),
		}))
		reminderFromDb, err := notesRepo.FindReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), r.ReminderID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminderFromDb).Should(
			notes.BeAReminder(t, createNoteResponse.NoteID, fakeUser.ID, startDate, startTime, timezone, interval, "", 0, time.Time{}, time.Now(), time.Now()))

		t.Run("Cannot schedule more than one reminder to the same note", func(t *testing.T) {
			_, err = app.ScheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), ScheduleReminderRequest{
				StartDate: startDate.Format(time.DateOnly),
				StartTime: startTime,
				Timezone:  timezone,
				Interval:  string(interval),
			})

			g.Expect(err).Should(
				MatchError(notes.ErrOnlyOneReminderAllowed))
		})

		t.Run("Delete a note and remove a reminder automatically successfully", func(t *testing.T) {
			err = app.DeleteNote(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String())

			g.Expect(err).Should(
				Not(HaveOccurred()))
		})
	})

	t.Run("Reschedule a reminder successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})
		createReminderResponse, err := app.ScheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), ScheduleReminderRequest{
			StartDate: "2023-05-24",
			StartTime: "20:19",
			Timezone:  "UTC",
			Interval:  "Monthly",
		})
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "20:33"
		interval := notes.Daily
		endsAfterN := uint(5)

		r, err := app.RescheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), createReminderResponse.ReminderID.String(), RescheduleReminderRequest{
			StartDate:  startDate.Format(time.DateOnly),
			StartTime:  startTime,
			Timezone:   timezone,
			Interval:   string(interval),
			EndsAfterN: endsAfterN,
		})

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(r).Should(gstruct.MatchAllFields(gstruct.Fields{
			"ReminderID": Not(Equal(uuid.Nil)),
			"CreatedAt":  BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt":  BeTemporally("~", time.Now(), time.Second),
		}))
		reminderFromDb, err := notesRepo.FindReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), createReminderResponse.ReminderID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminderFromDb).Should(
			notes.BeAReminder(t, createNoteResponse.NoteID, fakeUser.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
	})

	t.Run("Delete a reminder successfully", func(t *testing.T) {
		t.Parallel()

		fakeUser := notes.FakeUser(t)
		err := notesRepo.CreateUser(ctx, fakeUser)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		title := "test title"
		description := "test description"
		createNoteResponse, err := app.CreateNote(ctx, fakeUser.ID.String(), CreateNoteRequest{
			Title:       title,
			Description: description,
		})
		createReminderResponse, err := app.ScheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), ScheduleReminderRequest{
			StartDate: "2023-05-24",
			StartTime: "20:19",
			Timezone:  "UTC",
			Interval:  "Monthly",
		})
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = app.DeleteReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), createReminderResponse.ReminderID.String())

		g.Expect(err).Should(
			Not(HaveOccurred()))
		_, err = notesRepo.FindReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), createReminderResponse.ReminderID.String())
		g.Expect(err).Should(
			MatchError(notes.ErrReminderNotFound))
	})
}

func initializeApplication(_ *testing.T) (
	notes.Repository,
	CommandApplication,
	error,
) {
	db, err := repositories.NewGormDBFromEnv()
	if err != nil {
		return nil, CommandApplication{}, err
	}
	db = db.Debug()
	notesRepo := repositories.NewNotesRepository(db)
	app := NewCommandApplication(notesRepo, cron.NewLocalCron("/tmp/note-taking-tests"))

	return notesRepo, app, nil
}
