package cache

import (
	"context"

	"github.com/alejandro-cardenas-g/social/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStorage() Storage {
	return Storage{
		Users: &UsersMockStore{},
	}
}

type UsersMockStore struct {
	mock.Mock
}

func (s *UsersMockStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := s.Called(userID)
	return nil, args.Error(1)
}
func (s *UsersMockStore) Set(ctx context.Context, user *store.User) error {
	args := s.Called(user)
	return args.Error(0)
}
