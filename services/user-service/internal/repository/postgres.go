package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	userv1 "github.com/skillofide/proto/user/v1"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// VerifyUser verifies credentials against the database.
func (r *UserRepository) VerifyUser(ctx context.Context, email, password string) (*userv1.VerifyUserResponse, error) {
	var id, name, dbPassword, role string
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, name, password, role
		FROM   users
		WHERE  email = $1
	`, email).Scan(&id, &name, &dbPassword, &role)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	// Plain-text check as currently done by api-gateway
	if password != dbPassword {
		return nil, fmt.Errorf("invalid password")
	}

	return &userv1.VerifyUserResponse{
		Id:    id,
		Email: email,
		Name:  name,
		Role:  role,
	}, nil
}

// CreateOrUpdateUser inserts or updates a user record.
func (r *UserRepository) CreateOrUpdateUser(ctx context.Context, email, name, password, role string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (email, name, password, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) 
		DO UPDATE SET name = EXCLUDED.name, password = EXCLUDED.password, role = EXCLUDED.role, updated_at = now();
	`, email, name, password, role)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

// EnsureUsersTable creates the users table if missing and seeds the default admin user.
func (r *UserRepository) EnsureUsersTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email      TEXT NOT NULL UNIQUE,
			name       TEXT NOT NULL,
			password   TEXT NOT NULL,
			role       TEXT NOT NULL DEFAULT 'student',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	// Seed default user
	_, err = r.pool.Exec(ctx, `
		INSERT INTO users (email, name, password, role)
		VALUES ('admin@skillofied.com', 'Admin User', 'skillofied123', 'admin')
		ON CONFLICT (email) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("seed default user: %w", err)
	}

	return nil
}
