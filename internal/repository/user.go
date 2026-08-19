package repository

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2,$3) RETURNING id, created_at `,
		user.Username, user.Email, user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r UserRepository) GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	user := &models.User{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, created_at, username, email, password_hash, karma,is_verified
		FROM users WHERE email = $1 OR username = $1`,
		identifier,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Karma,
		&user.IsVerified,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r UserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, created_at, username, email, password_hash, karma,is_verified
		FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Karma,
		&user.IsVerified,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r UserRepository) MarkVerified(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE users SET is_verified = true WHERE id = $1`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("mark verified user: %w", err)
	}
	return nil
}
