package models

import (
	"github.com/google/uuid"
	"time"
)

type Vote struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	PostID    uuid.UUID `json:"post_id"`
	VoteType  int       `json:"vote_type"`
	CreatedAt time.Time `json:"created_at"`
}
