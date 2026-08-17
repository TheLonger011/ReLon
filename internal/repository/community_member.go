package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityMemberRepository struct {
	pool *pgxpool.Pool
}

func NewCommunityMemberRepository(pool *pgxpool.Pool) *CommunityMemberRepository {
	return &CommunityMemberRepository{pool: pool}
}

func (r CommunityMemberRepository) AddMember(ctx context.Context, communityID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO community_members(community_id, user_id) VALUES ($1, $2)`,
		communityID, userID,
	)

	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (r CommunityMemberRepository) RemoveMember(ctx context.Context, communityID, userID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM community_members WHERE community_id = $1 AND user_id = $2`,
		communityID, userID,
	)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member not found or access denied")
	}
	return nil
}

func (r CommunityMemberRepository) IsMember(ctx context.Context, communityID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM community_members WHERE community_id = $1 AND user_id = $2)`,
		communityID, userID,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check if community member: %w", err)
	}
	return exists, nil
}
