package repository

import (
	"context"
	"fmt"
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
		INSERT INTO comments(author_id, post_id, content) 
		VALUES ($1, $2, $3)
		RETURNING id,created_at,updated_at`,
		comment.AuthorID, comment.PostID, comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r CommentRepository) GetCommentByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	comments := []models.Comment{}

	rows, err := r.pool.Query(ctx, `
		SELECT id,post_id,author_id,content,created_at,updated_at 
		FROM comments WHERE post_id = $1 
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		postID, limit, offset,
	)
	if err != nil {
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return comments, err
		}

		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return comments, err
	}
	return comments, nil
}

func (r CommentRepository) DeleteComment(ctx context.Context, commentID, authorID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM comments WHERE id = $1 AND author_id = $2`,
		commentID, authorID,
	)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete comment: no rows affected")
	}
	return nil
}

func (r CommentRepository) UpdateComment(ctx context.Context, commentID, authorID uuid.UUID, content string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 AND author_id = $3`,
		content, commentID, authorID,
	)
	if err != nil {
		return fmt.Errorf("update comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update comment: no rows affected")
	}
	return nil
}
