package user

import (
	"context"
	"log/slog"
	"sync"

	user_adapter "chattery/internal/adapter/postgres/user"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/sliceutil"
)

type UserStore struct {
	adapter         *user_adapter.Adapter
	usersByID       map[domain.UserID]*domain.User
	usersByUsername map[domain.Username]*domain.User
	users           []*domain.User
	mu              sync.RWMutex
}

func New(adapter *user_adapter.Adapter) *UserStore {
	return &UserStore{
		adapter:         adapter,
		usersByID:       make(map[domain.UserID]*domain.User),
		usersByUsername: make(map[domain.Username]*domain.User),
		users:           make([]*domain.User, 0),
	}
}

func (*UserStore) Name() string {
	return "user_store"
}

func (s *UserStore) Sync(ctx context.Context) error {
	users, err := s.adapter.List(ctx)
	if err != nil {
		return err
	}

	slog.Info("[user_store] update", slog.Int("len", len(users)))

	groupedByID := sliceutil.SliceToMap(users, groupUserByID)
	groupedByUsername := sliceutil.SliceToMap(users, groupUserByUsername)

	s.mu.Lock()
	s.usersByID = groupedByID
	s.usersByUsername = groupedByUsername
	s.users = users
	s.mu.Unlock()

	return nil
}

func (s *UserStore) GetByID(id domain.UserID) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByID[id]
	if !ok {
		return nil, errutil.E().Kind(errutil.NotFound).Messagef("user id='%d' not found", id)
	}
	return user, nil
}

func (s *UserStore) GetByUsername(username domain.Username) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.usersByUsername[username]
	if !ok {
		return nil, errutil.E().
			Kind(errutil.NotFound).
			Messagef("user username='%s' not found", username)
	}
	return user, nil
}

func (s *UserStore) ListByID() map[domain.UserID]*domain.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.usersByID
}

func (s *UserStore) List() []*domain.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.users
}

func groupUserByID(user *domain.User) (domain.UserID, *domain.User) {
	return user.ID, user
}

func groupUserByUsername(user *domain.User) (domain.Username, *domain.User) {
	return user.Username, user
}
