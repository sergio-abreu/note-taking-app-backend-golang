package notes

import (
	"context"
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
	usersRepo, notesRepo, _, app, err := initializeApplication(t)
	g.Expect(err).Should(
		Not(HaveOccurred()))

	t.Run("Create a note successfully", func(t *testing.T) {
		fakeUser := notes.FakeUser(t)
		err := usersRepo.CreateUser(ctx, fakeUser)
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
			"NoteID": Not(Equal(uuid.Nil)),
		}))
		noteFromDB, err := notesRepo.FindNote(ctx, fakeUser.ID.String(), r.NoteID.String())
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(noteFromDB).Should(
			notes.BeANote(t, title, description, false, fakeUser.ID, time.Now(), time.Time{}))
	})
}

func initializeApplication(_ *testing.T) (
	notes.UsersRepository,
	notes.NotesRepository,
	notes.RemindersRepository,
	Application,
	error,
) {
	db, err := repositories.NewGormDBFromEnv()
	if err != nil {
		return nil, nil, nil, Application{}, err
	}
	db = db.Debug()
	usersRepo := repositories.NewUsersRepository(db)
	notesRepo := repositories.NewNotesRepository(db)
	var remindersRepo notes.RemindersRepository
	app := NewApplication(
		usersRepo,
		notesRepo,
		remindersRepo,
	)

	return usersRepo, notesRepo, remindersRepo, app, nil
}
