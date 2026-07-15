package repository

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
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
		INSERT INTO users (username, email, password_hash) VALUES ($1, $2,$3) RETURNING id, created_at `,
		user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r UserRepository) GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error) {
	user := &models.User{}
	err := r.pool.QueryRow(ctx, `
		SELECT id, created_at, username, email, password_hash, karma
		FROM users WHERE email = $1 OR username = $1`, identifier).Scan(&user.ID, &user.CreatedAt, &user.Username, &user.Email, &user.PasswordHash, &user.Karma)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil

}
