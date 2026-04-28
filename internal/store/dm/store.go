package dm

import (
	"context"
	"log/slog"
	"sync"

	dm_adapter "chattery/internal/adapter/postgres/dm"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/sliceutil"
)

type DMStore struct {
	adapter *dm_adapter.Adapter
	dmsByID map[domain.DMID]*domain.DM
	dms     []*domain.DM
	m       sync.RWMutex
}

func New(adapter *dm_adapter.Adapter) *DMStore {
	return &DMStore{
		adapter: adapter,
		dmsByID: make(map[domain.DMID]*domain.DM),
		dms:     make([]*domain.DM, 0),
	}
}

func (*DMStore) Name() string {
	return "dm_store"
}

func (s *DMStore) Sync(ctx context.Context) error {
	dms, err := s.adapter.ListDMs(ctx)
	if err != nil {
		return err
	}

	slog.Info("[dm_store] update", slog.Int("len", len(dms)))

	groupedByID := sliceutil.SliceToMap(dms, groupDMByID)

	s.m.Lock()
	s.dms = dms
	s.dmsByID = groupedByID
	s.m.Unlock()

	return nil
}

func (s *DMStore) GetByID(id domain.DMID) (*domain.DM, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	dm, ok := s.dmsByID[id]
	if !ok {
		return nil, errutil.E().Kind(errutil.NotFound).Messagef("dm id='%d' not found", id.I64())
	}
	return dm, nil
}

func (s *DMStore) List() []*domain.DM {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.dms
}

func groupDMByID(dm *domain.DM) (domain.DMID, *domain.DM) {
	return dm.ID, dm
}
