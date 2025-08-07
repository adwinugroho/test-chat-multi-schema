package service

import (
	"context"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/model"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
)

type tenantService struct {
	tenantRepository domain.TenantRepository
}

func NewTenantService(tenantRepository domain.TenantRepository) domain.TenantService {
	return &tenantService{tenantRepository: tenantRepository}
}

func (s *tenantService) NewTenant(ctx context.Context, tenant *domain.Tenant) error {
	err := s.tenantRepository.Create(ctx, tenant)
	if err != nil {
		logger.LogError("Error while create tenant: " + err.Error())
		return model.NewError(model.ErrorGeneral, "Internal server error")
	}
	return nil
}

func (s *tenantService) RemoveTenantByID(ctx context.Context, id string) error {
	err := s.tenantRepository.Delete(ctx, id)
	if err != nil {
		logger.LogError("Error while delete tenant: " + err.Error())
		return model.NewError(model.ErrorGeneral, "Internal server error")
	}
	return nil
}

func (s *tenantService) CreateTenantPartition(ctx context.Context, tenantID string) error {
	err := s.tenantRepository.CreateTenantPartition(ctx, tenantID)
	if err != nil {
		logger.LogError("Error while create tenant partition: " + err.Error())
		return model.NewError(model.ErrorGeneral, "Internal server error")
	}
	return nil
}
