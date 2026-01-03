package signaling

import (
	"chattery/internal/domain"
	"context"
	"sync"

	"github.com/coder/websocket"
)

type chatService interface {
	SendMessage(ctx context.Context, chat domain.ChatID, message *domain.Message)
	Subscribe(ctx context.Context, chat domain.ChatID, dst chan<- domain.Event)
}

type Service struct {
	mSubs *sync.RWMutex
	subs  map[domain.Username]*Subscriber
}

func New(chatSvc chatService) *Service {
	return &Service{
		subs: make(map[domain.Username]*Subscriber),
	}
}

func (s *Service) Subscribe(ctx context.Context, user domain.Username, ws *websocket.Conn) (context.Context, *Subscriber) {
	ctx, cancel := context.WithCancel(ctx)

	sub := &Subscriber{
		user:   user,
		ws:     ws,
		cancel: cancel,
	}

	s.mSubs.Lock()
	s.subs[user] = sub
	s.mSubs.Unlock()

	return ctx, sub
}

func (s *Service) Unsubscribe(sub *Subscriber) {
	s.mSubs.Lock()
	delete(s.subs, sub.user)
	s.mSubs.Unlock()
}
