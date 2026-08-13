package repository

import (
	"context"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{pool: pool}
}

func (r CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	err := r.pool.QueryRow(ctx, `
	INSERT INTO comments(author_id, post_id, content) VALUES ($1, $2, $3)
	RETURNING id,created_at,updated_at`, comment.AuthorID, comment.PostID, comment.Content).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r CommentRepository) GetCommentByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	var comments []models.Comment

	rows, err := r.pool.Query(ctx, `
		SELECT id,post_id,author_id,content,created_at,updated_at FROM comments WHERE post_id = $1 
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)
	if err != nil {
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return comments, err
	}
	return comments, nil

}
