package cache

import (
	"context"
	"gopher_social/internal/store"
)

func NewMockStore() *Storage {
	return &Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Get(context.Context, int64) (*store.User, error) {
	return nil, nil
}

func (m *MockUserStore) Set(context.Context, *store.User) error {
	return nil
}

func (m *MockUserStore) Delete(context.Context, int64) error {
	return nil
}
