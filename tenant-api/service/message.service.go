package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/google/uuid"
)

type messageService struct {
	publisher  domain.PublisherService
	repository domain.MessageRepository
}

func NewMessageService(pub domain.PublisherService, repo domain.MessageRepository) domain.MessageService {
	return &messageService{
		publisher:  pub,
		repository: repo,
	}
}

func (s *messageService) PublishMessage(ctx context.Context, tenantID string, req *model.PublishRequest) error {
	jsonBytes, err := json.Marshal(req.Content)
	if err != nil {
		return err
	}
	message := &domain.Message{
		MessageID: uuid.NewString(),
		Payload:   jsonBytes,
		TenantID:  tenantID,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.publisher.Publish(ctx, tenantID, body)
}

func (s *messageService) GetMessages(ctx context.Context, qParam map[string]string) ([]domain.Message, error) {
	return s.repository.GetMessages(ctx, qParam)
}
