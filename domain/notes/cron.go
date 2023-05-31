package notes

import "context"

type Cron interface {
	CreateCron(ctx context.Context, reminder Reminder) error
	DeleteCron(ctx context.Context, reminder Reminder) error
}
