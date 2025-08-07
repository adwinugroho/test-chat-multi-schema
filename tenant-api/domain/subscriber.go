package domain

import "context"

type SubscribeService interface {
	ConsumeTenantQueue(ctx context.Context, tenantID string, stopChan <-chan struct{})
}
