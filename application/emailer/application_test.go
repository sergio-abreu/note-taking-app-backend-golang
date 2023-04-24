package emailer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/emailer"
	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/repositories"
)

func TestApplication(t *testing.T) {
	g := NewWithT(t)
	ctx := context.Background()

	t.Run("Send reminder email", func(t *testing.T) {
		t.Parallel()

		ctrl, notesRepo, mailer, app, err := initializeApplication(t)
		defer ctrl.Finish()

		g.Expect(err).Should(
			Not(HaveOccurred()))
		user := notes.FakeUser(t)
		err = notesRepo.CreateUser(ctx, user)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		note, err := user.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = notesRepo.CreateNote(ctx, note)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "0 0 1 * *", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = notesRepo.ScheduleReminder(ctx, reminder)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		mailer.EXPECT().
			SendNoteReminder(ctx, user.Email, notes.Wrap(notes.BeANote(t, note.Title, note.Description, note.Completed, note.UserID, note.CreatedAt, note.UpdatedAt))).
			Return(nil)

		err = app.SendReminderEmail(ctx, SendReminderEmailRequest{
			ReminderID: reminder.ID.String(),
			NoteID:     note.ID.String(),
			UserID:     user.ID.String(),
		})

		g.Expect(err).Should(
			Not(HaveOccurred()))
	})
}

func initializeApplication(t *testing.T) (
	*gomock.Controller,
	notes.NotesRepository,
	*emailer.MockEmailer,
	Application,
	error,
) {
	db, err := repositories.NewGormDBFromEnv()
	if err != nil {
		return nil, nil, nil, Application{}, err
	}
	db = db.Debug()
	ctrl := gomock.NewController(t)
	notesRepo := repositories.NewNotesRepository(db)
	mailer := emailer.NewMockEmailer(ctrl)
	app := NewApplication(
		notesRepo,
		mailer,
	)

	return ctrl, notesRepo, mailer, app, nil
}
