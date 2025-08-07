package service

import (
	"context"
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type publisherService struct {
	conn *amqp.Connection
}

func NewPublisherService(conn *amqp.Connection) domain.PublisherService {
	return &publisherService{conn: conn}
}

func (p *publisherService) Publish(ctx context.Context, tenantID string, body []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
