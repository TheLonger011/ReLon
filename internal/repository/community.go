package repository

import (
	"context"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRepository struct {
	pool *pgxpool.Pool
}

func NewCommunityRepository(pool *pgxpool.Pool) *CommunityRepository {
	return &CommunityRepository{pool: pool}
}

func (r CommunityRepository) CreateCommunity(ctx context.Context, community *models.Community) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO communities (owner_id,name,description,is_private) 
		values ($1,$2,$3,$4)
		RETURNING id, created_at, updated_at`,
		community.OwnerID, community.Name, community.Description, community.IsPrivate,
	).Scan(&community.ID, &community.CreatedAt, &community.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create community: %w", err)
	}
	return nil
}

func (r CommunityRepository) GetCommunityByID(ctx context.Context, communityID uuid.UUID) (*models.Community, error) {
	community := &models.Community{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, owner_id, name, description, is_private,created_at,updated_at 
		FROM communities WHERE id = $1`,
		communityID,
	).Scan(
		&community.ID,
		&community.OwnerID,
		&community.Name,
		&community.Description,
		&community.IsPrivate,
		&community.CreatedAt,
		&community.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get community: %w", err)
	}
	return community, nil
}

func (r CommunityRepository) GetCommunities(ctx context.Context, limit, offset int) ([]models.Community, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, owner_id, name, description,is_private,created_at,updated_at
		FROM communities ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get communities: %w", err)
	}

	defer rows.Close()

	communities := []models.Community{}

	for rows.Next() {
		var community models.Community

		err := rows.Scan(
			&community.ID,
			&community.OwnerID,
			&community.Name,
			&community.Description,
			&community.IsPrivate,
			&community.CreatedAt,
			&community.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("get community: %w", err)
		}
		communities = append(communities, community)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get communities: %w", err)
	}
	return communities, nil
}

func (r CommunityRepository) DeleteCommunity(ctx context.Context, communityID, ownerID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM communities WHERE id = $1 AND owner_id = $2`,
		communityID, ownerID,
	)
	if err != nil {
		return fmt.Errorf("delete community: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("community not found or access denied")
	}
	return nil
}

func (r CommunityRepository) UpdateCommunity(ctx context.Context, communityID, ownerID uuid.UUID, name, description string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE communities SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND owner_id = $4`,
		name, description, communityID, ownerID,
	)
	if err != nil {
		return fmt.Errorf("update community: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("community not found or already deleted")
	}
	return nil
}
