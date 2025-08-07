package service

import (
	"context"
	"encoding/json"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
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
		logger.LogError("Error while marshall content message:" + err.Error())
		return model.NewError(model.ErrorGeneral, "Internal server error")
	}

	err = s.publisher.Publish(ctx, tenantID, jsonBytes)
	if err != nil {
		logger.LogError("Error while publish message:" + err.Error())
		return model.NewError(model.ErrorGeneral, "Internal server error")
	}

	return nil
}

func (s *messageService) GetMessages(ctx context.Context, tenantID string, qParam map[string]string) ([]model.ListMessagesResponse, string, error) {
	messages, nextCursor, err := s.repository.GetMessages(ctx, tenantID, qParam)
	if err != nil {
		return nil, "", err
	}

	if len(messages) == 0 {
		logger.LogInfo("no rows")
		return nil, "", nil
	}

	var results = make([]model.ListMessagesResponse, 0)
	for _, m := range messages {
		var r = model.ListMessagesResponse{
			MessageID: m.MessageID,
			Content:   string(m.Payload),
			CreatedAt: m.CreatedAt,
		}
		results = append(results, r)
	}

	return results, nextCursor, nil
}
