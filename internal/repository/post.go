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
		RETURNING id, created_at, updated_at, likes_count, dislikes_count`,
		post.AuthorID, post.Title, post.Content,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.LikesCount, &post.DislikesCount)

	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	return nil
}

func (r PostRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	post := &models.Post{}

	err := r.pool.QueryRow(ctx, `
		SELECT id,author_id, title, content,created_at, likes_count, dislikes_count,updated_at
		FROM posts WHERE id = $1`, postID,
	).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.LikesCount, &post.DislikesCount, &post.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return post, nil

}

func (r PostRepository) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {

	rows, err := r.pool.Query(ctx, `
		SELECT id, author_id, title, content, created_at, likes_count, dislikes_count,updated_at
		FROM posts ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	defer rows.Close()

	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.LikesCount, &post.DislikesCount, &post.UpdatedAt)

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

func (r PostRepository) SearchPosts(ctx context.Context, query string, limit, offset int) ([]models.Post, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id,author_id, title, content, created_at, likes_count, dislikes_count,updated_at
		FROM posts WHERE title ILIKE '%' || $1 || '%'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		query, limit, offset,
	)

	if err != nil {
		return nil, fmt.Errorf("search posts: %w", err)
	}

	defer rows.Close()

	posts := []models.Post{}

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt, &post.LikesCount, &post.DislikesCount, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("search posts: %w", err)
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search posts: %w", err)
	}
	return posts, nil

}

func (r PostRepository) DeletePost(ctx context.Context, postID, authorID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM posts WHERE id = $1 AND author_id = $2`,
		postID, authorID,
	)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("post not found or access denied")
	}

	return nil
}

func (r PostRepository) UpdatePost(ctx context.Context, postID, authorID uuid.UUID, title, content string) error {
	tag, err := r.pool.Exec(ctx, `
    	UPDATE posts SET title = $1, content = $2, updated_at = NOW() WHERE id = $3 AND author_id = $4 `,
		title, content, postID, authorID,
	)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("post not found or access denied")
	}
	return nil
}
