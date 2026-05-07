package config

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (cfg Config) NewRabbitMQ() (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
