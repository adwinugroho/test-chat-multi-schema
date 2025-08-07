package repository

import (
	"context"
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/helper"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tenantPgRepo struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) domain.TenantRepository {
	return &tenantPgRepo{db: db}
}

func (r *tenantPgRepo) Create(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		INSERT INTO tenants (tenant_id, tenant_name, user_id)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query,
		tenant.TenantID,
		tenant.TenantName,
		tenant.UserID,
	)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}

func (r *tenantPgRepo) CreateTenantPartition(ctx context.Context, tenantID string) error {
	tableName := helper.SanitizeTenantID(tenantID)
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS messages_%s
		PARTITION OF messages
		FOR VALUES IN ('%s')
	`, tableName, tenantID)

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}

func (r *tenantPgRepo) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM tenants WHERE tenant_id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}
