package user

import (
	"context"
	"log/slog"
	"slices"
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
	user, ok := s.usersByID[id]
	s.mu.RUnlock()
	if ok {
		return user, nil
	}

	user, err := s.adapter.UserByID(context.Background(), id)
	if err != nil {
		return nil, errutil.E(err).Debug("s.adapter.UserByID")
	}

	s.writeUser(user)

	return user, nil
}

func (s *UserStore) GetByUsername(username domain.Username) (*domain.User, error) {
	s.mu.RLock()
	user, ok := s.usersByUsername[username]
	s.mu.RUnlock()
	if ok {
		return user, nil
	}

	user, err := s.adapter.UserByUsername(context.Background(), username)
	if err != nil {
		return nil, errutil.E(err).Debug("s.adapter.UserByUsername")
	}

	s.writeUser(user)

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

func (s *UserStore) writeUser(user *domain.User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if old, ok := s.usersByID[user.ID]; ok && old.Username != user.Username {
		delete(s.usersByUsername, old.Username)
	}

	s.usersByID[user.ID] = user
	s.usersByUsername[user.Username] = user

	existingIndex := slices.IndexFunc(s.users, func(existing *domain.User) bool {
		return existing.ID == user.ID
	})

	if existingIndex >= 0 {
		s.users[existingIndex] = user
	} else {
		s.users = append(s.users, user)
	}
}
