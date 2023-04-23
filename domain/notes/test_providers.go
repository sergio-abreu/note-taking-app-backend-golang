package notes

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
)

func BeANote(
	_ *testing.T,
	title string,
	description string,
	completed bool,
	userID uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) types.GomegaMatcher {
	return gstruct.MatchAllFields(gstruct.Fields{
		"ID":          gomega.Not(gomega.Equal(uuid.Nil)),
		"Title":       gomega.Equal(title),
		"Description": gomega.Equal(description),
		"Completed":   gomega.Equal(completed),
		"UserID":      gomega.Equal(userID),
		"CreatedAt":   gomega.BeTemporally("~", createdAt, time.Second),
		"UpdatedAt":   gomega.BeTemporally("~", updatedAt, time.Second),
	})
}

func BeAReminder(
	_ *testing.T,
	noteID uuid.UUID,
	userID uuid.UUID,
	cronExpression string,
	endsAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) types.GomegaMatcher {
	return gstruct.MatchAllFields(gstruct.Fields{
		"ID":             gomega.Not(gomega.Equal(uuid.Nil)),
		"NoteID":         gomega.Equal(noteID),
		"UserID":         gomega.Equal(userID),
		"CronExpression": gomega.Equal(cronExpression),
		"EndsAt":         gomega.BeTemporally("~", endsAt, time.Second),
		"CreatedAt":      gomega.BeTemporally("~", createdAt, time.Second),
		"UpdatedAt":      gomega.BeTemporally("~", updatedAt, time.Second),
	})
}

func FakeUser(_ *testing.T) User {
	fakeUser := gofakeit.Person()
	return User{
		ID:        uuid.FromStringOrNil(gofakeit.UUID()),
		Name:      fakeUser.FirstName,
		Email:     fakeUser.Contact.Email,
		CreatedAt: time.Now(),
	}
}

type matcher struct {
	types.GomegaMatcher
	x interface{}
}

func (m matcher) Matches(x interface{}) bool {
	m.x = x
	result, _ := m.Match(x)
	return result
}

func (m matcher) String() string {
	return m.FailureMessage(m.x)
}

func Wrap(gmatcher types.GomegaMatcher) gomock.Matcher {
	return matcher{gmatcher, nil}
}
