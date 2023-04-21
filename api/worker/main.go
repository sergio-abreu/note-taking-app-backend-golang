package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure/messaging"
)

type S struct {
	ReminderId string `json:"reminder_id"`
	NoteId     string `json:"note_id"`
	UserId     string `json:"user_id"`
}

func main() {
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
		var s S
		err = json.Unmarshal(msg.Body, &s)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		data, _ := json.Marshal(s)
		fmt.Println(string(data))
		err = msg.Ack(false)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	}
}
