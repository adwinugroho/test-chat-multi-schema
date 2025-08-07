package service

import (
	"context"
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type subscriberService struct {
	messageRepository domain.MessageRepository
	conn              *amqp.Connection
}

func NewListenSubscriber(m domain.MessageRepository, rmqConn *amqp.Connection) domain.SubscribeService {
	return &subscriberService{
		messageRepository: m,
		conn:              rmqConn,
	}
}

func (s *subscriberService) ConsumeTenantQueue(ctx context.Context, tenantID string, stopChan <-chan struct{}) {
	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)

	ch, err := s.conn.Channel()
	if err != nil {
		logger.LogWithFields(logrus.Fields{
			"tenant_id": tenantID,
			"error":     fmt.Sprintf("failed to open channel: %v", err),
		}, "Error open channel RMQ")
		return
	}
	defer func() {
		_, err := ch.QueueDelete(queueName, false, false, false)
		if err != nil {
			logger.LogWithFields(logrus.Fields{
				"tenant_id": tenantID,
				"error":     fmt.Sprintf("failed to remove queue: %v", err),
			}, "Error remove queue")
		}
		ch.Close()
	}()

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		logger.LogWithFields(logrus.Fields{
			"tenant_id": tenantID,
			"error":     fmt.Sprintf("failed to declare queue: %v", err),
		}, "Error declare queue RMQ")
		return
	}

	msgs, err := ch.Consume(
		queueName,
		"consumer_"+tenantID, // consumer tag
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.LogWithFields(logrus.Fields{
			"tenant_id":  tenantID,
			"queue_name": queueName,
			"error":      fmt.Sprintf("failed to consume: %v", err),
		}, "Error channel consume")
		return
	}

	logger.LogInfo("Awaiting to consume...")

	for {
		select {
		case <-stopChan:
			logger.LogInfo(fmt.Sprintf("stopping on signal:%s", queueName))
			return
		case msg, ok := <-msgs:
			if !ok {
				logger.LogInfo(fmt.Sprintf("delivery message is closed:%s", queueName))
				return
			}

			logger.LogInfo(fmt.Sprintf("received message:%s", string(msg.Body)))
			newPayload := &domain.Message{
				MessageID: uuid.NewString(),
				TenantID:  tenantID,
				Payload:   msg.Body,
			}

			err := s.messageRepository.SaveMessage(ctx, newPayload)
			if err != nil {
				logger.LogError("Error while save message:" + err.Error())
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}
}
