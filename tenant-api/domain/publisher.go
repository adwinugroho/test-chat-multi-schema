package domain

import "context"

type PublisherService interface {
	Publish(ctx context.Context, tenantID string, body []byte) error
}
