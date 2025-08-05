package domain

type Tenant struct {
	ID         int    `json:"id"`
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
}

type TenantRepositoryInterface interface{}

type TenantServiceInterface interface{}
