package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type VoteInterface interface {
	Vote(ctx context.Context, userID, postID uuid.UUID, voteType int) error
}

type VoteService struct {
	repo VoteInterface
}

var ErrInvalidVoteType = errors.New("invalid vote type")

func NewVoteService(repo VoteInterface) *VoteService {
	return &VoteService{repo: repo}
}

func (s *VoteService) Vote(ctx context.Context, userID, postID uuid.UUID, voteType int) error {
	if voteType != 1 && voteType != -1 {
		return ErrInvalidVoteType
	}
	return s.repo.Vote(ctx, userID, postID, voteType)
}
