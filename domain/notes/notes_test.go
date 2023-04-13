package notes

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gofrs/uuid"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

func TestNote_CreateNote(t *testing.T) {
	g := NewWithT(t)

	t.Run("Create note successfully", func(t *testing.T) {
		user := FakeUser(t)
		title := "test title"
		description := "test description"

		note, err := user.CreateNote(title, description)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note).Should(
			BeANote(title, description, false, user.ID, time.Now(), time.Time{}))
	})

	t.Run("Don't create note when title is empty", func(t *testing.T) {
		user := FakeUser(t)
		title := ""
		description := "test description"

		_, err := user.CreateNote(title, description)

		g.Expect(err).Should(
			MatchError(ErrEmptyTitle))
	})
}

func TestNote_EditNote(t *testing.T) {
	g := NewWithT(t)

	t.Run("Edit note successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		newTitle := "new title test"
		newDescription := "new description test"

		err = user.EditNote(&note, newTitle, newDescription)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note).Should(
			BeANote(newTitle, newDescription, false, user.ID, time.Now(), time.Now()))
	})

	t.Run("Don't edit note when title is empty", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		newTitle := ""
		newDescription := "new description test"

		err = user.EditNote(&note, newTitle, newDescription)

		g.Expect(err).Should(
			MatchError(ErrEmptyTitle))
	})

	t.Run("Don't edit note when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		newTitle := "new title test"
		newDescription := "new description test"

		err = user1.EditNote(&noteFromUser2, newTitle, newDescription)

		g.Expect(err).Should(
			MatchError(ErrNoteDoesntBelongToThisUser))
	})
}

func TestNote_MarkNoteAsComplete(t *testing.T) {
	g := NewWithT(t)

	t.Run("Complete note successfully", func(t *testing.T) {
		user := FakeUser(t)
		title := "title test"
		description := "description test"
		note, err := user.CreateNote(title, description)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.MarkNoteAsCompleted(&note)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note).Should(
			BeANote(title, description, true, user.ID, time.Now(), time.Now()))
	})

	t.Run("Don't complete a note when it's already completed", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = user.MarkNoteAsCompleted(&note)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.MarkNoteAsCompleted(&note)

		g.Expect(err).Should(
			MatchError(ErrNotIsAlreadyCompleted))
	})

	t.Run("Don't complete a note when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user1.MarkNoteAsCompleted(&noteFromUser2)

		g.Expect(err).Should(
			MatchError(ErrNoteDoesntBelongToThisUser))
	})
}

func TestNote_MarkNoteAsInProgress(t *testing.T) {
	g := NewWithT(t)

	t.Run("Mark note as in progress successfully", func(t *testing.T) {
		user := FakeUser(t)
		title := "title test"
		description := "description test"
		note, err := user.CreateNote(title, description)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = user.MarkNoteAsCompleted(&note)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.MarkNoteAsInProgress(&note)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note).Should(
			BeANote(title, description, false, user.ID, time.Now(), time.Now()))
	})

	t.Run("Don't mark note as in progress when it's already in progress", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.MarkNoteAsInProgress(&note)

		g.Expect(err).Should(
			MatchError(ErrNotIsAlreadyInProgress))
	})

	t.Run("Don't mark note as in progress when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = user2.MarkNoteAsCompleted(&noteFromUser2)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user1.MarkNoteAsInProgress(&noteFromUser2)

		g.Expect(err).Should(
			MatchError(ErrNoteDoesntBelongToThisUser))
	})
}

func TestNote_CopyNote(t *testing.T) {
	g := NewWithT(t)

	t.Run("Copy an in progress note successfully", func(t *testing.T) {
		user := FakeUser(t)
		title := "title test"
		description := "description test"
		note1, err := user.CreateNote(title, description)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		note2, err := user.CopyNote(note1)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note2).Should(
			BeANote(title, description, false, user.ID, time.Now(), time.Time{}))
	})

	t.Run("Copy a completed note successfully", func(t *testing.T) {
		user := FakeUser(t)
		title := "title test"
		description := "description test"
		note1, err := user.CreateNote(title, description)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = user.MarkNoteAsCompleted(&note1)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		note2, err := user.CopyNote(note1)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note2).Should(
			BeANote(title, description, false, user.ID, time.Now(), time.Time{}))
	})

	t.Run("Don't copy a note when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user1.CopyNote(noteFromUser2)

		g.Expect(err).Should(
			MatchError(ErrNoteDoesntBelongToThisUser))
	})
}

func TestNote_DeleteNote(t *testing.T) {
	g := NewWithT(t)

	t.Run("Delete note successfully", func(t *testing.T) {
		user := FakeUser(t)
		title := "title test"
		description := "description test"
		note, err := user.CreateNote(title, description)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.DeleteNote(note)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(note).Should(
			BeANote(title, description, false, user.ID, time.Now(), time.Time{}))
	})

	t.Run("Don't delete note when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user1.DeleteNote(noteFromUser2)

		g.Expect(err).Should(
			MatchError(ErrNoteDoesntBelongToThisUser))
	})
}

func BeANote(
	title string,
	description string,
	completed bool,
	userID uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) types.GomegaMatcher {
	return MatchAllFields(Fields{
		"ID":          Not(Equal(uuid.Nil)),
		"Title":       Equal(title),
		"Description": Equal(description),
		"Completed":   Equal(completed),
		"UserID":      Equal(userID),
		"CreatedAt":   BeTemporally("~", createdAt, time.Second),
		"UpdatedAt":   BeTemporally("~", updatedAt, time.Second),
	})
}

func FakeUser(_ *testing.T) User {
	fakeUser := gofakeit.Person()
	return User{
		ID:    uuid.FromStringOrNil(gofakeit.UUID()),
		Name:  fakeUser.FirstName,
		Email: fakeUser.Contact.Email,
	}
}
