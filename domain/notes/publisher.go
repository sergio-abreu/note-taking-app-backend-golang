package notes

import "context"

type Publisher interface {
	PublishReminderScheduled(ctx context.Context, reminder Reminder) error
	PublishReminderRescheduled(ctx context.Context, reminder Reminder) error
	PublishReminderDeleted(ctx context.Context, reminder Reminder) error
}
