package domain

import "context"

type Tenant struct {
	ID         int    `json:"id"`
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
}

type TenantService interface {
	NewTenant(ctx context.Context, tenant *Tenant) error
	RemoveTenantByID(ctx context.Context, id string) error
}
