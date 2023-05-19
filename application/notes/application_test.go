package notes

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
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
		cronExpression := "33 20 19 * *"

		r, err := app.ScheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), ScheduleReminderRequest{
			CronExpression: cronExpression,
			EndsAt:         "",
			Repeats:        0,
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
			notes.BeAReminder(t, createNoteResponse.NoteID, fakeUser.ID, cronExpression, time.Time{}, time.Now(), time.Now()))

		t.Run("Cannot schedule more than one reminder to the same note", func(t *testing.T) {
			_, err = app.ScheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), ScheduleReminderRequest{
				CronExpression: cronExpression,
				EndsAt:         "",
				Repeats:        0,
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
			CronExpression: "33 20 19 * *",
			EndsAt:         "",
			Repeats:        0,
		})
		g.Expect(err).Should(
			Not(HaveOccurred()))
		repeats := uint(5)
		refDate := time.Now().AddDate(0, 0, -1)
		cronExpression := fmt.Sprintf("28 21 %d * *", refDate.Day())
		endsAtByRepetition := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 21, 28, 0, 0, refDate.Location()).
			AddDate(0, int(repeats), 0)

		r, err := app.RescheduleReminder(ctx, fakeUser.ID.String(), createNoteResponse.NoteID.String(), createReminderResponse.ReminderID.String(), RescheduleReminderRequest{
			CronExpression: cronExpression,
			EndsAt:         "",
			Repeats:        repeats,
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
			notes.BeAReminder(t, createNoteResponse.NoteID, fakeUser.ID, cronExpression, endsAtByRepetition, time.Now(), time.Now()))
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
			CronExpression: "33 20 19 * *",
			EndsAt:         "",
			Repeats:        0,
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
	app := NewCommandApplication(notesRepo)

	return notesRepo, app, nil
}
