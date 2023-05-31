package messaging

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewPublisher(ch *amqp.Channel) (Publisher, error) {
	err := ch.ExchangeDeclare("reminders", amqp.ExchangeFanout, true, false, false, true, nil)
	if err != nil {
		return Publisher{}, err
	}
	err = ch.Confirm(false)
	if err != nil {
		return Publisher{}, err
	}
	return Publisher{ch: ch}, nil
}

type Publisher struct {
	ch *amqp.Channel
}

func (p Publisher) PublishReminderScheduled(ctx context.Context, reminder notes.Reminder) error {
	body, err := json.Marshal(reminder)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, "reminders", "reminders.scheduled", true, true, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    reminder.ID.String(),
		Timestamp:    time.Now(),
		Type:         "ReminderScheduled",
		UserId:       reminder.UserID.String(),
		AppId:        "note-taking",
		Body:         body,
	})
}

func (p Publisher) PublishReminderRescheduled(ctx context.Context, reminder notes.Reminder) error {
	body, err := json.Marshal(reminder)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, "reminders", "reminders.rescheduled", true, true, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    reminder.ID.String(),
		Timestamp:    time.Now(),
		Type:         "ReminderRescheduled",
		UserId:       reminder.UserID.String(),
		AppId:        "note-taking",
		Body:         body,
	})
}

func (p Publisher) PublishReminderDeleted(ctx context.Context, reminder notes.Reminder) error {
	body, err := json.Marshal(reminder)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, "reminders", "reminders.deleted", true, true, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    reminder.ID.String(),
		Timestamp:    time.Now(),
		Type:         "ReminderDeleted",
		UserId:       reminder.UserID.String(),
		AppId:        "note-taking",
		Body:         body,
	})
}
