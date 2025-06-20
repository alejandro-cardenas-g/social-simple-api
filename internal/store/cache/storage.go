package cache

import (
	"context"

	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Users interface {
		Get(ctx context.Context, userID int64) (*store.User, error)
		Set(ctx context.Context, user *store.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{rdb: rdb},
	}
}
