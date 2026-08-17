package service

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
)

type CommunityJoinRequestInterface interface {
	CreateJoinRequest(ctx context.Context, communityID, userID uuid.UUID) error
	GetJoinRequestByID(ctx context.Context, requestID uuid.UUID) (*models.CommunityJoinRequest, error)
	UpdateJoinRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error
	GetPendingRequestByCommunity(ctx context.Context, communityID uuid.UUID) ([]models.CommunityJoinRequest, error)
}

type CommunityJoinRequestService struct {
	requestRepo   CommunityJoinRequestInterface
	memberRepo    CommunityMemberInterface
	communityRepo CommunityReaderInterface
}

func NewCommunityJoinRequestService(requestRepo CommunityJoinRequestInterface, memberRepo CommunityMemberInterface, communityRepo CommunityReaderInterface) *CommunityJoinRequestService {
	return &CommunityJoinRequestService{requestRepo: requestRepo, memberRepo: memberRepo, communityRepo: communityRepo}
}

func (s *CommunityJoinRequestService) JoinCommunity(ctx context.Context, communityID, userID uuid.UUID) error {
	community, err := s.communityRepo.GetCommunityByID(ctx, communityID)
	if err != nil {
		return fmt.Errorf("get community: %w", err)
	}
	user, err := s.memberRepo.IsMember(ctx, communityID, userID)
	if err != nil {
		return fmt.Errorf("check if user is member: %w", err)
	}
	if user {
		return fmt.Errorf("user is already a member")
	}
	if !community.IsPrivate {
		return s.memberRepo.AddMember(ctx, communityID, userID)
	}
	return s.requestRepo.CreateJoinRequest(ctx, communityID, userID)
}

func (s *CommunityJoinRequestService) ApproveJoinRequest(ctx context.Context, requestID, ownerID uuid.UUID) error {
	request, err := s.requestRepo.GetJoinRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get join request: %w", err)
	}
	community, err := s.communityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		return fmt.Errorf("get community: %w", err)
	}
	if community.OwnerID != ownerID {
		return fmt.Errorf("not author")
	}
	if err := s.memberRepo.AddMember(ctx, community.ID, request.UserID); err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	if err := s.requestRepo.UpdateJoinRequestStatus(ctx, requestID, "approved"); err != nil {
		return fmt.Errorf("update join request status: %w", err)
	}
	return nil
}

func (s *CommunityJoinRequestService) RejectJoinRequest(ctx context.Context, requestID, ownerID uuid.UUID) error {
	request, err := s.requestRepo.GetJoinRequestByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get join request: %w", err)
	}
	community, err := s.communityRepo.GetCommunityByID(ctx, request.CommunityID)
	if err != nil {
		return fmt.Errorf("get community: %w", err)
	}
	if community.OwnerID != ownerID {
		return fmt.Errorf("not author")
	}
	if err := s.requestRepo.UpdateJoinRequestStatus(ctx, requestID, "rejected"); err != nil {
		return fmt.Errorf("update join request status: %w", err)
	}
	return nil
}

func (s *CommunityJoinRequestService) GetPendingRequests(ctx context.Context, communityID, ownerID uuid.UUID) ([]models.CommunityJoinRequest, error) {
	community, err := s.communityRepo.GetCommunityByID(ctx, communityID)
	if err != nil {
		return nil, fmt.Errorf("get community: %w", err)
	}
	if community.OwnerID != ownerID {
		return nil, fmt.Errorf("not author")
	}
	return s.requestRepo.GetPendingRequestByCommunity(ctx, communityID)
}
