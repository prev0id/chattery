package server

import (
	"context"
	"log/slog"
	"sync"

	server_adapter "chattery/internal/adapter/postgres/server"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/sliceutil"
)

type ServerStore struct {
	adapter     *server_adapter.Adapter
	serversByID map[domain.ServerID]*domain.Server
	servers     []*domain.Server
	m           sync.RWMutex
}

func New(adapter *server_adapter.Adapter) *ServerStore {
	return &ServerStore{
		adapter:     adapter,
		serversByID: make(map[domain.ServerID]*domain.Server),
		servers:     make([]*domain.Server, 0),
	}
}

func (*ServerStore) Name() string {
	return "server_store"
}

func (s *ServerStore) Sync(ctx context.Context) error {
	servers, err := s.adapter.ListServers(ctx)
	if err != nil {
		return err
	}

	slog.Info("[server_store] update", slog.Int("len", len(servers)))

	groupedByID := sliceutil.SliceToMap(servers, groupServerByID)

	s.m.Lock()
	s.servers = servers
	s.serversByID = groupedByID
	s.m.Unlock()

	return nil
}

func (s *ServerStore) GetByID(id domain.ServerID) (*domain.Server, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	server, ok := s.serversByID[id]
	if !ok {
		return nil, errutil.E().Kind(errutil.NotFound).Messagef("server id='%d' not found", id.I64())
	}
	return server, nil
}

func (s *ServerStore) List() []*domain.Server {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.servers
}

func groupServerByID(server *domain.Server) (domain.ServerID, *domain.Server) {
	return server.ID, server
}
