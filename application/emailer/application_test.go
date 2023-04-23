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

		ctrl, usersRepo, notesRepo, remindersRepo, mailer, app, err := initializeApplication(t)
		defer ctrl.Finish()

		g.Expect(err).Should(
			Not(HaveOccurred()))
		user := notes.FakeUser(t)
		err = usersRepo.CreateUser(ctx, user)
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
		err = remindersRepo.ScheduleReminder(ctx, reminder)
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
	notes.UsersRepository,
	notes.NotesRepository,
	notes.RemindersRepository,
	*emailer.MockEmailer,
	Application,
	error,
) {
	db, err := repositories.NewGormDBFromEnv()
	if err != nil {
		return nil, nil, nil, nil, nil, Application{}, err
	}
	db = db.Debug()
	ctrl := gomock.NewController(t)
	usersRepo := repositories.NewUsersRepository(db)
	notesRepo := repositories.NewNotesRepository(db)
	remindersRepo := repositories.NewRemindersRepository(db)
	mailer := emailer.NewMockEmailer(ctrl)
	app := NewApplication(
		usersRepo,
		notesRepo,
		mailer,
	)

	return ctrl, usersRepo, notesRepo, remindersRepo, mailer, app, nil
}

// {8b79ba3e-fe67-4274-8699-83265c1f601f test title test description false daab9ab9-8162-41f6-9743-765fcc2bddb9 2023-04-23 16:19:08.235967 -0300 -03 0001-01-01 00:00:00 +0000 UTC}
// {8b79ba3e-fe67-4274-8699-83265c1f601f test title test description false daab9ab9-8162-41f6-9743-765fcc2bddb9 2023-04-23 16:19:08.235967 -0300 -03 0001-01-01 00:00:00 +0000 UTC}

/*
2023-04-23 16:24:41.061040113 -0300 m=+0.022221915
2023-04-23 16:24:41.06104 -0300

0001-01-01 00:00:00 +0000
0001-01-01 00:00:00 +0000
*/
