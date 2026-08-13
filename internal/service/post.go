package service

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
)

type PostInterface interface {
	CreatePost(ctx context.Context, post *models.Post) error
	GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error)
}

type PostService struct {
	repo PostInterface
}

func NewPostService(repo PostInterface) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(
	ctx context.Context,
	authorID uuid.UUID,
	title, content string) (*models.Post, error) {

	post := models.Post{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
	}
	err := s.repo.CreatePost(ctx, &post)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return &post, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	return s.repo.GetPostByID(ctx, postID)
}

func (s *PostService) GetPosts(ctx context.Context, limit, offset int) ([]models.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetPosts(ctx, limit, offset)
}
