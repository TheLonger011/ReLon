package models

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID            uuid.UUID  `json:"id"`
	AuthorID      uuid.UUID  `json:"author_id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	CreatedAt     time.Time  `json:"created_at"`
	LikesCount    int        `json:"likes_count"`
	DislikesCount int        `json:"dislikes_count"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CommunityID   *uuid.UUID `json:"community_id,omitempty"`
}
