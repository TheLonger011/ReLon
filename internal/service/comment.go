package service

import (
	"context"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
)

type CommentInterface interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	GetCommentByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]models.Comment, error)
}

type CommentService struct {
	repo CommentInterface
}

func NewCommentService(repo CommentInterface) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(ctx context.Context, authorID, postID uuid.UUID, content string) (*models.Comment, error) {
	comment := &models.Comment{
		AuthorID: authorID,
		PostID:   postID,
		Content:  content,
	}
	err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]models.Comment, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetCommentByPostID(ctx, postID, limit, offset)
}
