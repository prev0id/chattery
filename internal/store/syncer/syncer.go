package syncer

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
)

type Store interface {
	Name() string
	Sync(context.Context) error
}

type Syncer[T Store] struct {
	store   T
	stopCh  chan struct{}
	timeout time.Duration
}

func Start[T Store](timeout time.Duration, store T) error {
	s := &Syncer[T]{
		store:   store,
		timeout: timeout,
		stopCh:  make(chan struct{}),
	}

	if err := s.store.Sync(context.Background()); err != nil {
		return errors.E(err).Debug("s.store.Sync")
	}

	go s.run()

	return nil
}

func (s *Syncer[T]) run() {
	ticker := time.NewTicker(s.timeout)
	defer ticker.Stop()

	slog.Info("syncer started", slog.String("name", s.store.Name()))
	defer slog.Warn("syncer stopped", slog.String("name", s.store.Name()))

	for {
		select {
		case <-s.stopCh:
			slog.Info("syncer started", slog.String("name", s.store.Name()))
			return
		case <-ticker.C:
			s.sync()
		}
	}
}

func (s *Syncer[T]) sync() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error(
				"syncer recovered from panic",
				slog.String("name", s.store.Name()),
				slog.Any("panic", r),
				slog.String("stack", string(debug.Stack())),
			)
		}
	}()

	if err := s.store.Sync(context.Background()); err != nil {
		logger.Error(err, "syncer.sync")
	}
}
