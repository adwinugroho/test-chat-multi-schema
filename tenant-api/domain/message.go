package domain

import (
	"context"
	"time"

	"github.com/adwinugroho/test-chat-multi-schema/model"
)

type Message struct {
	MessageID string    `db:"message_id"`
	TenantID  string    `db:"tenant_id"`
	Payload   []byte    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
}

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *Message) error
	GetMessages(ctx context.Context, qParam map[string]string) ([]Message, error)
}

type MessageService interface {
	PublishMessage(ctx context.Context, tenantID string, req *model.PublishRequest) error
	GetMessages(ctx context.Context, qParam map[string]string) ([]Message, error)
}
