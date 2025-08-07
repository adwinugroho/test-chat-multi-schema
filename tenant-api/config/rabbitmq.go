package config

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EnvRabbitMQConfig struct {
	RabbitMQ RabbitMQConn `mapstructure:"rabbitmq"`
}

type RabbitMQConn struct {
	URL string `mapstructure:"url"`
}

var (
	RabbitMQConfig EnvRabbitMQConfig
	RabbitConn     *amqp.Connection
)

func InitRabbitMQConnection(conn string) error {
	rmqConnection, err := amqp.Dial(conn)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	RabbitConn = rmqConnection
	return nil
}
