package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sergio-abreu/note-taking-app-backend-golang/application/emailer"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/emails"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/messaging"
	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/repositories"
)

func main() {
	db, err := repositories.NewGormDBFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	db = db.Debug()

	mailer, err := emails.NewEMailerFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	app := emailer.NewApplication(repositories.NewNotesRepository(db), mailer)

	conn, err := messaging.NewRabbitmqFromEnv()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	err = ch.ExchangeDeclare("reminders", amqp.ExchangeFanout, true, false, false, true, nil)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	queue, err := ch.QueueDeclare("send_reminder_email", true, false, false, false, nil)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	err = ch.QueueBind(queue.Name, "", "reminders", false, nil)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	delivery, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for msg := range delivery {
		var r emailer.SendReminderEmailRequest
		err = json.Unmarshal(msg.Body, &r)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = app.SendReminderEmail(context.Background(), r)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = msg.Ack(false)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}
