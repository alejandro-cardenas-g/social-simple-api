package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (s *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}
func (s *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	return &User{}, nil
}

func (s *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}
func (s *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
func (s *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
func (s *MockUserStore) GetByEmail(ctx context.Context, Email string) (*User, error) {
	return &User{}, nil
}
