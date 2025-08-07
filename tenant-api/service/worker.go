package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
)

type TenantConsumer struct {
	stopChan chan struct{}
	doneChan chan struct{}
}

type TenantManager struct {
	mu        sync.Mutex
	consumers map[string]*TenantConsumer
	subscribe domain.SubscribeService
}

func NewTenantManager(subs domain.SubscribeService) *TenantManager {
	return &TenantManager{
		consumers: make(map[string]*TenantConsumer),
		subscribe: subs,
	}
}

func (tm *TenantManager) StartConsumer(ctx context.Context, tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.consumers[tenantID]; exists {
		return fmt.Errorf("consumer for tenant %s already running", tenantID)
	}

	stopChan := make(chan struct{})
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)
		tm.subscribe.ConsumeTenantQueue(ctx, tenantID, stopChan)
	}()

	tm.consumers[tenantID] = &TenantConsumer{
		stopChan: stopChan,
		doneChan: doneChan,
	}

	return nil
}

func (tm *TenantManager) StopConsumer(tenantID string) error {
	tm.mu.Lock()
	consumer, exists := tm.consumers[tenantID]
	if !exists {
		tm.mu.Unlock()
		return fmt.Errorf("no consumer for tenant %s", tenantID)
	}
	delete(tm.consumers, tenantID)
	tm.mu.Unlock()

	close(consumer.stopChan)
	<-consumer.doneChan
	return nil
}

// optional to kill all running queue
func (tm *TenantManager) StopAllConsumers() {
	tm.mu.Lock()
	tenantIDs := make([]string, 0, len(tm.consumers))
	for tenantID := range tm.consumers {
		tenantIDs = append(tenantIDs, tenantID)
	}
	tm.mu.Unlock()

	for _, tenantID := range tenantIDs {
		_ = tm.StopConsumer(tenantID)
	}
}
