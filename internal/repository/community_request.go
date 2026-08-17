package repository

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRequestRepository struct {
	pool *pgxpool.Pool
}

func NewCommunityRequestRepository(pool *pgxpool.Pool) *CommunityRequestRepository {
	return &CommunityRequestRepository{pool: pool}
}

func (r CommunityRequestRepository) CreateJoinRequest(ctx context.Context, communityID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO community_join_requests(community_id, user_id) 
		VALUES ($1, $2)`,
		communityID, userID,
	)
	if err != nil {
		return fmt.Errorf("create join request: %w", err)
	}
	return nil
}

func (r CommunityRequestRepository) GetJoinRequestByID(ctx context.Context, requestID uuid.UUID) (*models.CommunityJoinRequest, error) {
	joinRequest := &models.CommunityJoinRequest{}
	err := r.pool.QueryRow(ctx, `
		SELECT id,community_id,user_id,status,created_at
		FROM community_join_requests WHERE id = $1`,
		requestID,
	).Scan(
		&joinRequest.ID,
		&joinRequest.CommunityID,
		&joinRequest.UserID,
		&joinRequest.Status,
		&joinRequest.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get join request: %w", err)
	}
	return joinRequest, nil
}

func (r CommunityRequestRepository) UpdateJoinRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	rows, err := r.pool.Exec(ctx, `
		UPDATE community_join_requests SET status = $1 WHERE id = $2`,
		status, requestID,
	)
	if err != nil {
		return fmt.Errorf("update join request status: %w", err)
	}
	if rows.RowsAffected() == 0 {
		return fmt.Errorf("join request not found or access denied")
	}
	return nil
}

func (r CommunityRequestRepository) GetPendingRequestByCommunity(ctx context.Context, communityID uuid.UUID) ([]models.CommunityJoinRequest, error) {
	rows, err := r.pool.Query(ctx, `
			SELECT id,community_id,user_id,status,created_at
			FROM community_join_requests WHERE community_id = $1 AND status = 'pending'`,
		communityID,
	)
	if err != nil {
		return nil, fmt.Errorf("get pending request: %w", err)
	}
	defer rows.Close()
	requests := []models.CommunityJoinRequest{}
	for rows.Next() {
		var community models.CommunityJoinRequest
		err := rows.Scan(
			&community.ID,
			&community.CommunityID,
			&community.UserID,
			&community.Status,
			&community.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("get pending request: %w", err)
		}
		requests = append(requests, community)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get pending request: %w", err)
	}
	return requests, nil
}
