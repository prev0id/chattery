package chat

import (
	"context"
	"slices"

	"chattery/internal/domain"
	"chattery/internal/utils/logger"
)

func (s *Service) Register(sub domain.Subscriber) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	s.subsByUserID[sub.GetUserID()] = append(s.subsByUserID[sub.GetUserID()], sub)
}

func (s *Service) Unregister(sub domain.Subscriber) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	s.deleteFromSubsBySession(sub)
	s.deleteFromSubsByUserID(sub)
}

func (s *Service) StartListeningToChat(ctx context.Context, sub domain.Subscriber, chat domain.ChatID) {
	ctx = s.newSubscriptionContexWithCancel(ctx, sub)

	messages := make(chan *domain.Message)
	s.pubsub.Subscribe(ctx, chat, messages)

	go func() {
		for message := range messages {
			s.broadcastMessage(ctx, chat, message)
		}
	}()
}

func (s *Service) StopListeningToChat(sub domain.Subscriber) {
	s.cancelChatSubscription(sub)
}

func (s *Service) deleteFromSubsBySession(target domain.Subscriber) {
	delete(s.chatSubsBySession, target.GetSession())
}

func (s *Service) deleteFromSubsByUserID(target domain.Subscriber) {
	for userID, subs := range s.subsByUserID {
		newSubs := slices.DeleteFunc(subs, func(sub domain.Subscriber) bool {
			return sub.GetSession() == target.GetSession()
		})

		if len(newSubs) == 0 {
			delete(s.subsByUserID, userID)
		} else {
			s.subsByUserID[userID] = newSubs
		}
	}
}

func (s *Service) broadcastMessage(ctx context.Context, chat domain.ChatID, message *domain.Message) {
	participants, err := s.db.ListParticipants(ctx, chat)
	if err != nil {
		logger.ErrorCtx(ctx, err, "s.db.ListParticipants")
		return
	}

	for _, participant := range participants {
		s.subsMutex.RLock()
		defer s.subsMutex.RUnlock()
		subs, exist := s.subsByUserID[participant.UserID]

		if !exist {
			continue
		}

		for _, sub := range subs {
			if err := sub.WriteMessage(ctx, message); err != nil {
				logger.ErrorCtx(ctx, err, "send message error")
			}
		}
	}
}

func (s *Service) newSubscriptionContexWithCancel(ctx context.Context, sub domain.Subscriber) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	currentCancel, exists := s.chatSubsBySession[sub.GetSession()]
	if exists { // отменяем прошлое соединение если оно существует
		currentCancel()
	}

	s.chatSubsBySession[sub.GetSession()] = cancel

	return ctx
}

func (s *Service) cancelChatSubscription(sub domain.Subscriber) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()

	currentCancel, exists := s.chatSubsBySession[sub.GetSession()]
	if exists { // отменяем прошлое соединение если оно существует
		currentCancel()
	}
}
