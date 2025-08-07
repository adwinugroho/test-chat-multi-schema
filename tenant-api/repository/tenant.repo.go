package repository

import (
	"context"
	"fmt"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
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
		INSERT INTO tenant (id, tenant_id, tenant_name)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query,
		tenant.TenantID,
		tenant.TenantName,
	)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}

func (r *tenantPgRepo) CreateTenantPartition(ctx context.Context, tenantID string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS messages_%s
		PARTITION OF messages
		FOR VALUES IN ('%s')
	`, tenantID, tenantID)

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}

func (r *tenantPgRepo) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM tenant WHERE tenant_id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}
