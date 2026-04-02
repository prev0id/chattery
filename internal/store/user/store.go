package user_store

import (
	"context"
	"sync"

	user_adapter "chattery/internal/adapter/postgres/user"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/sliceutil"
)

type UserStore struct {
	adapter   *user_adapter.Adapter
	mu        sync.RWMutex
	usersByID map[domain.UserID]*domain.User
}

func New(adapter *user_adapter.Adapter) *UserStore {
	return &UserStore{
		adapter:   adapter,
		usersByID: make(map[domain.UserID]*domain.User),
	}
}

func (s *UserStore) Name() string {
	return "user_store"
}

func (s *UserStore) Sync(ctx context.Context) error {
	users, err := s.adapter.List(ctx)
	if err != nil {
		return err
	}

	groupedByID := sliceutil.SliceToMap(users, func(user *domain.User) (domain.UserID, *domain.User) {
		return user.ID, user
	})

	s.mu.Lock()
	s.usersByID = groupedByID
	s.mu.Unlock()

	return nil
}

func (s *UserStore) GetByID(id domain.UserID) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[id]
	if !ok {
		return nil, errors.E().Kind(errors.NotFound).Messagef("user id='%d' not found", id)
	}
	return user, nil
}

func (s *UserStore) ListByIDs(ids ...domain.UserID) map[domain.UserID]*domain.User {
	result := make(map[domain.UserID]*domain.User, len(ids))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range ids {
		if user, ok := s.usersByID[id]; ok {
			result[id] = user
		}
	}
	return result
}
