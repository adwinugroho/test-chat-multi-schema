package config

import (
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type EnvRabbitMQConfig struct {
	RabbitMQ RabbitMQConn `mapstructure:"rabbitmq"`
}

type RabbitMQConn struct {
	URL string `mapstructure:"url"`
}

var (
	RabbitMQConfig EnvRabbitMQConfig
)

func InitRabbitMQConnection(conn string) (*amqp.Connection, error) {
	logger.LogWithFields(logrus.Fields{
		"info": "Connecting to rabbitMQ",
		"url":  conn,
	}, "info connecting to message broker")
	rmqConnection, err := amqp.Dial(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return rmqConnection, nil
}
