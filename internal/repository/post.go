package repository

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

func NewPostRepository(pool *pgxpool.Pool) *PostRepository {
	return &PostRepository{pool: pool}
}

func (r PostRepository) CreatePost(ctx context.Context, post *models.Post) error {

	err := r.pool.QueryRow(ctx, `
		INSERT INTO posts(author_id, title, content) VALUES ($1, $2, $3) 
		RETURNING id, created_at, likes_count, dislikes_count`,
		post.AuthorID, post.Title, post.Content,
	).Scan(&post.ID, &post.CreatedAt, &post.LikesCount, &post.DislikesCount)

	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	return nil
}

func (r PostRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post := &models.Post{}

	err := r.pool.QueryRow(ctx, `
		SELECT id,author_id, title, content,created_at, likes_count, dislikes_count
		FROM posts WHERE id = $1`, postID,
	).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.LikesCount, &post.DislikesCount)

	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return post, nil

}

func (r PostRepository) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {

	rows, err := r.pool.Query(ctx, `
		SELECT id, author_id, title, content, created_at, likes_count, dislikes_count
		FROM posts ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.LikesCount, &post.DislikesCount)

		if err != nil {
			return nil, fmt.Errorf("get posts: %w", err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	return posts, nil
}
