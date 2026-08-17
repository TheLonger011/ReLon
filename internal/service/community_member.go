package service

import (
	"context"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
)

type CommunityMemberInterface interface {
	AddMember(ctx context.Context, communityID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, communityID, userID uuid.UUID) error
	IsMember(ctx context.Context, communityID, userID uuid.UUID) (bool, error)
}

type CommunityReaderInterface interface {
	GetCommunityByID(ctx context.Context, communityID uuid.UUID) (*models.Community, error)
}

type CommunityMemberService struct {
	memberRepo    CommunityMemberInterface
	communityRepo CommunityReaderInterface
}

func NewCommunityMemberService(memberRepo CommunityMemberInterface, communityRepo CommunityReaderInterface) *CommunityMemberService {
	return &CommunityMemberService{memberRepo: memberRepo, communityRepo: communityRepo}
}

func (s *CommunityMemberService) LeaveCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	return s.memberRepo.RemoveMember(ctx, communityID, userID)
}
