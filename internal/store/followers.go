package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type FollowersStore struct {
	db *sql.DB
}

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

func (s *FollowersStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	query := `
		INSERT INTO followers (follower_id, user_id) 
		VALUES ($1, $2);
	`

	_, err := s.db.ExecContext(ctx, query, followerID, userID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}

	return err
}

func (s *FollowersStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)

	defer cancel()

	query := `
		DELETE FROM followers
		WHERE follower_id = $1 AND user_id = $2 
	`

	_, err := s.db.ExecContext(ctx, query, followerID, userID)

	return err
}
