package notes

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
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
			BeANote(t, title, description, false, user.ID, time.Now(), time.Now()))
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
			BeANote(t, newTitle, newDescription, false, user.ID, time.Now(), time.Now()))
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
			BeANote(t, title, description, true, user.ID, time.Now(), time.Now()))
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
			BeANote(t, title, description, false, user.ID, time.Now(), time.Now()))
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
			BeANote(t, title, description, false, user.ID, time.Now(), time.Now()))
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
			BeANote(t, title, description, false, user.ID, time.Now(), time.Now()))
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
			BeANote(t, title, description, false, user.ID, time.Now(), time.Now()))
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

func TestNote_ScheduleAReminder(t *testing.T) {
	g := NewWithT(t)

	t.Run("Schedule a reminder to repeat daily forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Daily
		endsAfterN := uint(0)

		reminder, err := user.ScheduleAReminder(note, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("0 1 * * *"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Schedule a reminder to repeat weekly (Mon, Wed, Fri) forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "13:23"
		interval := Weekly
		weekDays := "1,3,5"
		endsAfterN := uint(0)

		reminder, err := user.ScheduleAReminder(note, startDate.Format(time.DateOnly), startTime, timezone, string(interval), weekDays, "", endsAfterN)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, weekDays, endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("23 13 * * 1,3,5"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Schedule a reminder to repeat monthly forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "15:19"
		interval := Monthly
		endsAfterN := uint(0)

		reminder, err := user.ScheduleAReminder(note, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal(fmt.Sprintf("19 15 %d * *", refDate.Day())))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Schedule a reminder to repeat yearly forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "23:11"
		interval := Yearly
		endsAfterN := uint(0)

		reminder, err := user.ScheduleAReminder(note, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal(fmt.Sprintf("11 23 %d %d *", refDate.Day(), refDate.Month())))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Schedule a reminder to repeat 3 times successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now().UTC().AddDate(0, 0, -1)
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Daily
		endsAfterN := uint(3)
		endsAtByRepetition := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 1, 0, 0, 0, time.UTC).
			AddDate(0, 0, int(endsAfterN))

		reminder, err := user.ScheduleAReminder(note, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("0 1 * * *"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeTemporally("~", endsAtByRepetition, time.Second))
	})

	t.Run("Schedule a reminder to repeat until certain date successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "00:12"
		interval := Daily
		endsAt := startDate.AddDate(0, 1, 0)
		endsAfterN := uint(0)

		reminder, err := user.ScheduleAReminder(note, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", endsAt.Format(time.DateOnly), endsAfterN)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, endsAt, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("12 0 * * *"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeTemporally("~", endsAt, time.Second))
	})

	t.Run("Don't create a reminder when start date format is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "invalid-format", "", "", "", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidStartDate))
	})

	t.Run("Don't create a reminder when start time format is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "2023-05-24", "invalid-format", "", "", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidStartTime))
	})

	t.Run("Don't create a reminder when timezone is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "2023-05-24", "10:02", "invalid-format", "", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidTimezone))
	})

	t.Run("Don't create a reminder when interval is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "invalid-format", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidInterval))
	})

	t.Run("Don't create a reminder when week days is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Weekly", "invalid-format", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidWeekDays))
	})

	t.Run("Don't create a reminder when endsAt date format is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Daily", "", "invalid-format", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidEndsAt))
	})

	t.Run("Don't create a reminder when configured multiple termination", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Daily", "", "2023-05-25", 1)

		g.Expect(err).Should(
			MatchError(ErrCannotConfigureMultipleTermination))
	})

	t.Run("Don't create a reminder when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		_, err = user1.ScheduleAReminder(noteFromUser2, "2023-05-24", "10:02", "UTC", "Daily", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrNoteDoesntBelongToThisUser))
	})
}

func TestNote_RescheduleAReminder(t *testing.T) {
	g := NewWithT(t)

	t.Run("Reschedule a reminder to repeat daily forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Daily
		endsAfterN := uint(0)

		err = user.RescheduleAReminder(&reminder, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("0 1 * * *"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Reschedule a reminder to repeat weekly (Tue, Thu, Sat, Sun) forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Weekly
		weekDays := "2,4,6,7"
		endsAfterN := uint(0)

		err = user.RescheduleAReminder(&reminder, startDate.Format(time.DateOnly), startTime, timezone, string(interval), weekDays, "", endsAfterN)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, weekDays, endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("0 1 * * 2,4,6,7"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Reschedule a reminder to repeat monthly forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Daily", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Monthly
		endsAfterN := uint(0)

		err = user.RescheduleAReminder(&reminder, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal(fmt.Sprintf("0 1 %d * *", refDate.Day())))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Reschedule a reminder to repeat yearly forever successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Daily", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Yearly
		endsAfterN := uint(0)

		err = user.RescheduleAReminder(&reminder, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal(fmt.Sprintf("0 1 %d %d *", refDate.Day(), refDate.Month())))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeZero())
	})

	t.Run("Reschedule a reminder to repeat 3 times successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now().UTC().AddDate(0, 0, -1)
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		fmt.Println(startDate.AddDate(0, 0, 3))
		timezone := "UTC"
		startTime := "01:00"
		interval := Daily
		endsAfterN := uint(3)
		endsAtByRepetition := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 1, 0, 0, 0, refDate.Location()).
			AddDate(0, 0, int(endsAfterN))

		err = user.RescheduleAReminder(&reminder, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", "", endsAfterN)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, time.Time{}, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("0 1 * * *"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeTemporally("~", endsAtByRepetition, time.Second))
	})

	t.Run("Reschedule a reminder to repeat until certain date successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		refDate := time.Now()
		startDate := time.Date(refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, time.UTC)
		timezone := "UTC"
		startTime := "01:00"
		interval := Daily
		endsAt := startDate.AddDate(0, 3, 0)
		endsAfterN := uint(0)

		err = user.RescheduleAReminder(&reminder, startDate.Format(time.DateOnly), startTime, timezone, string(interval), "", endsAt.Format(time.DateOnly), endsAfterN)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		g.Expect(reminder).Should(
			BeAReminder(t, note.ID, user.ID, startDate, startTime, timezone, interval, "", endsAfterN, endsAt, time.Now(), time.Now()))
		cron := reminder.ParseCron()
		g.Expect(cron).Should(
			Equal("0 1 * * *"))
		g.Expect(reminder.ParseEndsAt(cron)).Should(
			BeTemporally("~", endsAt, time.Second))
	})

	t.Run("Don't reschedule a reminder when start date format is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "invalid-format", "", "", "", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidStartDate))
	})

	t.Run("Don't reschedule a reminder when start date format is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "2023-05-24", "invalid-format", "", "", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidStartTime))
	})

	t.Run("Don't reschedule a reminder when timezone is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "2023-05-24", "10:02", "invalid", "", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidTimezone))
	})

	t.Run("Don't reschedule a reminder when interval is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "2023-05-24", "10:02", "UTC", "invalid", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidInterval))
	})

	t.Run("Don't reschedule a reminder when week days is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "2023-05-24", "10:02", "UTC", "Weekly", "invalid", "", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidWeekDays))
	})

	t.Run("Don't reschedule a reminder when endsAt date format is invalid", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "2023-05-24", "10:02", "UTC", "Daily", "", "invalid", 0)

		g.Expect(err).Should(
			MatchError(ErrInvalidEndsAt))
	})

	t.Run("Don't reschedule a reminder when configured multiple termination", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.RescheduleAReminder(&reminder, "2023-05-24", "10:02", "UTC", "Daily", "", "invalid", 1)

		g.Expect(err).Should(
			MatchError(ErrCannotConfigureMultipleTermination))
	})

	t.Run("Don't reschedule a reminder when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminderFromUser2, err := user2.ScheduleAReminder(noteFromUser2, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user1.RescheduleAReminder(&reminderFromUser2, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)

		g.Expect(err).Should(
			MatchError(ErrReminderDoesntBelongToThisUser))
	})
}

func TestNote_DeleteReminder(t *testing.T) {
	g := NewWithT(t)

	t.Run("Delete reminder successfully", func(t *testing.T) {
		user := FakeUser(t)
		note, err := user.CreateNote("title test", "description test")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminder, err := user.ScheduleAReminder(note, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user.DeleteReminder(reminder)

		g.Expect(err).Should(
			Not(HaveOccurred()))
	})

	t.Run("Don't delete note when it doesn't belong to this user", func(t *testing.T) {
		user1 := FakeUser(t)
		user2 := FakeUser(t)
		noteFromUser2, err := user2.CreateNote("test title", "test description")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		reminderFromUser2, err := user2.ScheduleAReminder(noteFromUser2, "2023-05-24", "10:02", "UTC", "Monthly", "", "", 0)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = user1.DeleteReminder(reminderFromUser2)

		g.Expect(err).Should(
			MatchError(ErrReminderDoesntBelongToThisUser))
	})
}
