package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	pool *pgxpool.Pool
}

func NewVoteRepository(pool *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{pool: pool}
}

func (r VoteRepository) GetVote(ctx context.Context, userID, postID uuid.UUID) (*models.Vote, error) {
	vote := &models.Vote{}

	err := r.pool.QueryRow(ctx, `
		SELECT id,user_id,post_id, vote_type, created_at
		FROM votes WHERE user_id = $1 AND post_id = $2`,
		userID, postID,
	).Scan(&vote.ID, &vote.UserID, &vote.PostID, &vote.VoteType, &vote.CreatedAt)

	if err != nil {
		return nil, err
	}
	return vote, nil
}

func (r VoteRepository) Vote(ctx context.Context, userID, postID uuid.UUID, voteType int) error {
	vote := &models.Vote{}

	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		SELECT id,user_id,post_id, vote_type, created_at
		FROM votes WHERE user_id = $1 AND post_id = $2`, userID, postID,
	).Scan(&vote.ID, &vote.UserID, &vote.PostID, &vote.VoteType, &vote.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if err := updateCounter(ctx, tx, postID, voteType, 1); err != nil {
				return err
			}

			_, err = tx.Exec(ctx, `
				INSERT INTO votes(user_id, post_id, vote_type) 
				VALUES ($1, $2, $3)`,
				userID, postID, voteType,
			)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if vote.VoteType == voteType {
			_, err = tx.Exec(ctx,
				`DELETE FROM votes WHERE user_id = $1 AND post_id = $2`,
				userID, postID,
			)

			if err != nil {
				return err
			}

			if err := updateCounter(ctx, tx, postID, voteType, -1); err != nil {
				return err
			}
		} else {
			_, err = tx.Exec(ctx, `
				UPDATE votes SET vote_type = $1 WHERE user_id = $2 AND post_id = $3
				`, voteType, userID, postID,
			)

			if err != nil {
				return err
			}

			if err := updateCounter(ctx, tx, postID, vote.VoteType, -1); err != nil {
				return err
			}

			if err := updateCounter(ctx, tx, postID, voteType, 1); err != nil {
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func updateCounter(ctx context.Context, tx pgx.Tx, postID uuid.UUID, voteType int, delta int) error {
	column := "likes_count"

	if voteType != 1 {
		column = "dislikes_count"
	}

	_, err := tx.Exec(ctx, fmt.Sprintf(
		`UPDATE posts SET %s = %s + $1 WHERE id = $2`,
		column, column), delta, postID,
	)

	if err != nil {
		return err
	}

	return nil
}
