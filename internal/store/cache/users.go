package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/go-redis/redis/v8"
)

type UsersStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute * 5

func (s *UsersStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user:%v", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UsersStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user:%v", user.ID)

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, userData, UserExpTime).Err()
}
