package domain

import (
	"context"
	"time"
)

type Message struct {
	MessageID string    `db:"message_id"`
	TenantID  string    `db:"tenant_id"`
	Payload   []byte    `db:"payload"`
	CreatedAt time.Time `db:"created_at"`
}

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *Message) error
}
