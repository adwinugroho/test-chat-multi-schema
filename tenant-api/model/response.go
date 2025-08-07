package model

import "time"

// user
type AuthenticationResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

// tenant
type CreateTenantResponse struct {
	TenantID string `json:"tenant_id"`
	Message  string `json:"message"`
}

// messages
type ListMessagesResponse struct {
	MessageID string    `json:"message_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
