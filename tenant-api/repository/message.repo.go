package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/helper"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messagePgRepo struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &messagePgRepo{db: db}
}

func (r *messagePgRepo) GetMessages(ctx context.Context, tenantID string, qParam map[string]string) ([]domain.Message, string, error) {
	limit := qParam["limit"]
	limitInt, _ := strconv.Atoi(limit)

	cursor := qParam["cursor"]
	cursorInt64, _ := strconv.ParseInt(cursor, 10, 64)

	query := fmt.Sprintf(`
		SELECT message_id, payload, created_at
		FROM messages_%s
		WHERE ($1 = 0 OR EXTRACT(EPOCH FROM created_at) > $1)
		ORDER BY created_at ASC
		LIMIT $2;`, helper.SanitizeTenantID(tenantID))

	rows, err := r.db.Query(ctx, query, cursorInt64, limitInt+1)
	if err != nil {
		logger.LogError("Error while querying: " + err.Error())
		return nil, "", err
	}
	defer rows.Close()

	var results []domain.Message
	var nextCursor string
	for rows.Next() {
		var res domain.Message

		err := rows.Scan(
			&res.MessageID, &res.Payload, &res.CreatedAt,
		)
		if err != nil {
			logger.LogError("Error while scan row: " + err.Error())
			return nil, "", err
		}

		results = append(results, res)
	}

	if len(results) > limitInt {
		nextUnix := results[limitInt].CreatedAt.Unix()
		nextCursor = fmt.Sprintf("%d", nextUnix)
		results = results[:limitInt]
	}

	return results, nextCursor, nil
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
