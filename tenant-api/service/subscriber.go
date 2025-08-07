package service

import (
	"context"
	"fmt"
	"sync"

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

func (s *subscriberService) ConsumeTenantQueue(ctx context.Context, tenantID string, stopChan <-chan struct{}, workers int) {
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
		"consumer_"+tenantID,
		false, false, false, false, nil,
	)
	if err != nil {
		logger.LogWithFields(logrus.Fields{
			"tenant_id": tenantID,
			"error":     fmt.Sprintf("failed to consume queue: %v", err),
		}, "Error consume queue RMQ")
		return
	}

	logger.LogInfo(fmt.Sprintf("Starting %d workers for tenant %s", workers, tenantID))

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					logger.LogInfo(fmt.Sprintf("worker-%d stopping for tenant:%s", workerID, tenantID))
					return
				case msg, ok := <-msgs:
					if !ok {
						logger.LogInfo(fmt.Sprintf("channel closed for tenant:%s", tenantID))
						return
					}

					logger.LogInfo(fmt.Sprintf("worker-%d received message: %s", workerID, string(msg.Body)))

					newPayload := &domain.Message{
						MessageID: uuid.NewString(),
						TenantID:  tenantID,
						Payload:   msg.Body,
					}

					err := s.messageRepository.SaveMessage(ctx, newPayload)
					if err != nil {
						logger.LogError(fmt.Sprintf("worker-%d error saving message: %v", workerID, err))
						msg.Nack(false, true)
						continue
					}

					msg.Ack(false)
				}
			}
		}(i)
	}

	<-stopChan
	logger.LogInfo(fmt.Sprintf("stop signal received, waiting for workers of tenant %s to finish...", tenantID))
	wg.Wait()
}
