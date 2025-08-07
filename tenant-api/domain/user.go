package domain

import (
	"context"

	"github.com/adwinugroho/test-chat-multi-schema/model"
)

type User struct {
	UserID   string  `db:"user_id"`
	Name     string  `db:"name"`
	Email    string  `db:"email"`
	Password *string `db:"password"`
	Role     string  `db:"role"`
	TenantID *string `db:"tenant_id"` //NULLABLE
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	LoginUser(ctx context.Context, req model.LoginUserRequest) (*model.AuthenticationResponse, error)
}
