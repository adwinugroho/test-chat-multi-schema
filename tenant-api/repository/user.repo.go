package repository

import (
	"context"
	"errors"

	"github.com/adwinugroho/test-chat-multi-schema/domain"
	"github.com/adwinugroho/test-chat-multi-schema/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userPgRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userPgRepo{db: db}
}

func (r *userPgRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (user_id, email, password, name, role)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query,
		user.UserID,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
	)
	if err != nil {
		logger.LogError("Error querying: " + err.Error())
		return err
	}

	return nil
}

func (r *userPgRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT user_id, email, password, name, role, tenant_id FROM users WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.TenantID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.LogError("Error get user email: " + err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *userPgRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT user_id, email, password, name, role, tenant_id FROM users WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.TenantID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.LogError("Error user not found")
			return nil, nil
		}
		logger.LogError("Error get user ID: " + err.Error())
		return nil, err
	}

	return &user, nil
}
