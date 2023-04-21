package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/sergio-abreu/note-taking-app-backend-golang/infrastructure"
)

func NewRabbitmq(user, pass, host, port, vHost string) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vHost)
	return amqp.Dial(url)
}

func NewRabbitmqFromEnv() (*amqp.Connection, error) {
	user := infrastructure.GetEnvWithDefault("RABBITMQ_USER", "note-taking")
	pass := infrastructure.GetEnvWithDefault("RABBITMQ_PASSWORD", "note-taking")
	host := infrastructure.GetEnvWithDefault("RABBITMQ_HOST", "127.0.0.1")
	port := infrastructure.GetEnvWithDefault("RABBITMQ_PORT", "5672")
	vHost := infrastructure.GetEnvWithDefault("RABBITMQ_VHOST", "note-taking")
	return NewRabbitmq(user, pass, host, port, vHost)
}
