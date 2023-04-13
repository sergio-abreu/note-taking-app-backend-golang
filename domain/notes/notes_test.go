package notes

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gofrs/uuid"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
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
		g.Expect(note).Should(MatchAllFields(Fields{
			"ID":          Not(Equal(uuid.Nil)),
			"Title":       Equal(title),
			"Description": Equal(description),
			"UserID":      Equal(user.ID),
			"CreatedAt":   BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt":   BeZero(),
		}))
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
		g.Expect(note).Should(MatchAllFields(Fields{
			"ID":          Not(Equal(uuid.Nil)),
			"Title":       Equal(newTitle),
			"Description": Equal(newDescription),
			"UserID":      Equal(user.ID),
			"CreatedAt":   BeTemporally("~", time.Now(), time.Second),
			"UpdatedAt":   BeTemporally("~", time.Now(), time.Second),
		}))
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

	t.Run("Don't edit note when note doesn't belong to this user", func(t *testing.T) {
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

func FakeUser(_ *testing.T) User {
	fakeUser := gofakeit.Person()
	return User{
		ID:    uuid.FromStringOrNil(gofakeit.UUID()),
		Name:  fakeUser.FirstName,
		Email: fakeUser.Contact.Email,
	}
}
