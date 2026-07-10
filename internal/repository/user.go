package repository

import (
	"context"
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
