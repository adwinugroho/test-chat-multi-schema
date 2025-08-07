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

func (r *messagePgRepo) GetMessages(ctx context.Context, qParam map[string]string) ([]domain.Message, error) {
	limit := qParam["pageSize"]
	offset := qParam["offset"]

	query := `SELECT message_id, payload, created_at FROM messages ORDER BY created_at desc LIMIT $1 OFFSET $2;`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		logger.LogError("Error while querying: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []domain.Message
	for rows.Next() {
		var res domain.Message

		err := rows.Scan(
			&res.MessageID, &res.Payload, &res.CreatedAt,
		)
		if err != nil {
			logger.LogError("Error while scan row: " + err.Error())
			return nil, err
		}

		results = append(results, res)
	}

	return results, nil
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
