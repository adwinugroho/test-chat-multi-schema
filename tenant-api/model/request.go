package model

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type NewTenantRequest struct {
	TenantName string `json:"tenant_name" validate:"required"`
}
type PublishRequest struct {
	Content string `json:"content" validate:"required"`
}
