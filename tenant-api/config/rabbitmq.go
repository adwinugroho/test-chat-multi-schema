package config

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EnvRabbitMQConfig struct {
	RabbitMQ RabbitMQConn `mapstructure:"rabbitmq"`
}

type RabbitMQConn struct {
	URL string `mapstructure:"url"`
}

type Exchange struct {
	ExchangeKey string
	RoutingKey  string
	QueueName   string
}

var (
	RabbitConn     *amqp.Connection
	RabbitChann    *amqp.Channel
	RabbitMQConfig EnvRabbitMQConfig
	RabbitMQName   string
	RabbitMQClosed chan bool
)

func InitConnectionRabbitMQ(projectName, projectModule string) error {
	conn, err := amqp.Dial(RabbitMQConfig.RabbitMQ.URL)
	if err != nil {
		log.Printf("error when open to connection: %s, cause:%+v\n", RabbitMQConfig.RabbitMQ.URL, err)
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Println("error when open to channel connection: ", err)
		return err
	}

	topicName := fmt.Sprintf("%s_%s_rpc", projectName, projectModule)
	_, err = ch.QueueDeclare(
		topicName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		log.Printf("Error declaring queue, cause: %+v\n", err)
		return err
	}

	RabbitConn = conn
	RabbitChann = ch
	RabbitMQName = topicName

	log.Println("RabbitMQ successfully connected")
	return nil
}
