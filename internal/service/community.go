package service

import (
	"context"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
)

type CommunityInterface interface {
	CreateCommunity(ctx context.Context, community *models.Community) error
	GetCommunityByID(ctx context.Context, communityID uuid.UUID) (*models.Community, error)
	GetCommunities(ctx context.Context, limit, offset int) ([]models.Community, error)
	DeleteCommunity(ctx context.Context, communityID, ownerID uuid.UUID) error
	UpdateCommunity(ctx context.Context, communityID, ownerID uuid.UUID, name, description string) error
}

type CommunityService struct {
	repo CommunityInterface
}

func NewCommunityService(repo CommunityInterface) *CommunityService {
	return &CommunityService{repo: repo}
}

func (s *CommunityService) CreateCommunity(ctx context.Context, ownerID uuid.UUID, name, description string, isPrivate bool) (*models.Community, error) {
	community := models.Community{
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
	}
	if err := s.repo.CreateCommunity(ctx, &community); err != nil {
		return nil, err
	}

	return &community, nil
}

func (s *CommunityService) GetCommunityByID(ctx context.Context, communityID uuid.UUID) (*models.Community, error) {
	return s.repo.GetCommunityByID(ctx, communityID)
}

func (s *CommunityService) GetCommunities(ctx context.Context, limit, offset int) ([]models.Community, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetCommunities(ctx, limit, offset)
}

func (s *CommunityService) DeleteCommunity(ctx context.Context, communityID, ownerID uuid.UUID) error {
	return s.repo.DeleteCommunity(ctx, communityID, ownerID)
}

func (s *CommunityService) UpdateCommunity(ctx context.Context, communityID, ownerID uuid.UUID, name, description string) error {
	return s.repo.UpdateCommunity(ctx, communityID, ownerID, name, description)
}
