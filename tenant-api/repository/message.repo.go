package repository

import (
	"context"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messagePgRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &messagePgRepo{db: db}
}

func (r *messagePgRepo) SaveMessage(ctx context.Context, message *domain.Message) error {
	query := `
			INSERT INTO messages (message_id, tenant_id, payload, created_at)
			VALUES ($1, $2, $3, NOW())
		`
	_, err := r.db.Exec(ctx, query, message.MessageID, message.TenantID, message.Payload)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}
