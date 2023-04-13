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

func FakeUser(_ *testing.T) User {
	fakeUser := gofakeit.Person()
	return User{
		ID:    uuid.FromStringOrNil(gofakeit.UUID()),
		Name:  fakeUser.FirstName,
		Email: fakeUser.Contact.Email,
	}
}
