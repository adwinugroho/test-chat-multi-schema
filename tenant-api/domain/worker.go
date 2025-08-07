package domain

import "context"

type TenantManagerWorker interface {
	StartConsumer(ctx context.Context, tenantID string, workers int) error
	StopConsumer(tenantID string) error
	RestartConsumer(ctx context.Context, tenantID string, workers int) error
	StopAllConsumers()
}
