package domain

import "context"

type Tenant struct {
	TenantID   string `db:"tenant_id"`
	TenantName string `db:"tenant_name"`
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	CreateTenantPartition(ctx context.Context, tenantID string) error
	Delete(ctx context.Context, id string) error
}

type TenantService interface {
	CreateTenantPartition(ctx context.Context, tenantID string) error
	NewTenant(ctx context.Context, tenant *Tenant) error
	RemoveTenantByID(ctx context.Context, id string) error
}
