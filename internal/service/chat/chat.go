package chat

import (
	"context"
	"strings"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/sliceutil"
)

func (s *Service) JoinChat(ctx context.Context, user domain.UserID, chat domain.ChatID) error {
	participant := &domain.Participant{
		Role: domain.ChatRoleParticipant,
		User: user,
		Chat: chat,
	}

	if err := s.db.AddParticipant(ctx, participant); err != nil {
		return errors.E(err).Debug("s.db.AddParticipant")
	}
	return nil
}

func (s *Service) LeaveChat(ctx context.Context, user domain.UserID, chat domain.ChatID) error {
	if err := s.db.DeleteParticipant(ctx, user, chat); err != nil {
		return errors.E(err).Debug("s.db.DeleteParticipant")
	}
	return nil
}

func (s *Service) CreatePublicChat(ctx context.Context, user domain.UserID, name string) (domain.ChatID, error) {
	var (
		chatID domain.ChatID
		err    error
	)

	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		chat := &domain.Chat{
			Name: name,
			Type: domain.ChatTypePublic,
		}
		chatID, err = s.db.CreateChat(ctx, chat)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateChat")
		}

		moderator := &domain.Participant{
			Role: domain.ChatRoleOwner,
			User: user,
			Chat: chatID,
		}
		if err := s.db.AddParticipant(ctx, moderator); err != nil {
			return errors.E(err).Debug("s.db.AddParticipant")
		}
		return nil
	})

	return chatID, err
}

func (s *Service) CreatePrivateChat(ctx context.Context, users ...domain.UserID) (domain.ChatID, error) {
	var (
		chatID domain.ChatID
		err    error
	)
	err = s.transaction.InTransaction(ctx, func(ctx context.Context) error {
		chat := &domain.Chat{
			Type: domain.ChatTypePrivate,
		}
		chatID, err := s.db.CreateChat(ctx, chat)
		if err != nil {
			return errors.E(err).Debug("s.db.CreateChat")
		}

		for _, user := range users {
			if err := s.JoinChat(ctx, user, chatID); err != nil {
				return errors.E(err).Debug("s.db.JoinChat")
			}
		}
		return nil
	})

	return chatID, err
}

func (s *Service) SearchChats(ctx context.Context, query string) ([]*domain.Chat, error) {
	chats, err := s.db.Chats(ctx)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.Chats")
	}

	query = strings.ToLower(query)

	chats = sliceutil.Filter(chats, func(chat *domain.Chat) bool {
		return chatNameContainsString(chat, query)
	})

	return chats, nil
}

func (s *Service) DeleteChat(ctx context.Context, user domain.UserID, chat domain.ChatID) error {
	participants, err := s.db.ListParticipants(ctx, chat)
	if err != nil {
		return errors.E(err).Debug("s.db.ListParticipants")
	}

	owner, ok := sliceutil.Find(participants, func(participant *domain.Participant) bool {
		return participant.Role == domain.ChatRoleOwner
	})
	if !ok {
		return errors.E().Message("chat does not have owner").Kind(errors.NotFound)
	}
	if owner.User != user {
		return errors.E().Message("you must be chat owner to delete it").Kind(errors.Permission)
	}

	if err := s.db.DeleteChat(ctx, chat); err != nil {
		return errors.E(err).Debug("s.db.DeleteChat")
	}

	return nil
}

func (s *Service) UserChats(ctx context.Context, user domain.UserID) ([]*domain.Chat, error) {
	chats, err := s.db.UserChats(ctx, user)
	if err != nil {
		return nil, errors.E(err).Debug("s.db.UserChats")
	}
	return chats, nil
}

func chatNameContainsString(chat *domain.Chat, query string) bool {
	if chat.Type == domain.ChatTypePrivate {
		return false
	}
	chatName := strings.ToLower(chat.Name)
	return strings.Contains(chatName, query)
}
