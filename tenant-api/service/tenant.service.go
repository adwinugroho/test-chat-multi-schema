package service

import (
	"context"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
)

type tenantService struct {
	tenantRepository domain.TenantRepository
}

func NewTenantService(tenantRepository domain.TenantRepository) domain.TenantService {
	return &tenantService{tenantRepository: tenantRepository}
}

func (s *tenantService) NewTenant(ctx context.Context, tenant *domain.Tenant) error {
	return s.tenantRepository.Create(ctx, tenant)
}

func (s *tenantService) RemoveTenantByID(ctx context.Context, id string) error {
	return s.tenantRepository.Delete(ctx, id)
}
