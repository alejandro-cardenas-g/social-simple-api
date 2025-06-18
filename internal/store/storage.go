package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("RESOURCE NOT FOUND")
	ErrConflict          = errors.New("CONFLICT")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, postID int64) (*Post, error)
		UpdateByID(ctx context.Context, post *Post) error
		DeleteByID(ctx context.Context, postID int64) error
		GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
	}

	Users interface {
		Create(ctx context.Context, tx *sql.Tx, user *User) error
		GetByID(ctx context.Context, userID int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(ctx context.Context, token string) error
	}

	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID int64, userID int64) error
		Unfollow(ctx context.Context, followerID int64, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
	}
}

func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
