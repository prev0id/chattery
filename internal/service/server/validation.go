package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

func (s *Service) validateServerExists(ctx context.Context, serverID domain.ServerID) error {
	_, err := s.db.GetServer(ctx, serverID)
	if errors.Is(errors.NotFound, err) {
		return errors.E(err).Messagef("server id='%d' not found", serverID.I64())
	}

	if err != nil {
		return errors.E(err).Debug("v.db.GetServer")
	}

	return nil
}

func (s *Service) validateTopicExists(ctx context.Context, topicID domain.TopicID) error {
	_, err := s.db.GetTopic(ctx, topicID)
	return err
}

func (s *Service) validateParticipantExists(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	_, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errors.Is(errors.NotFound, err) {
		return errors.E(err).Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetServerParticipant")
	}
	return nil
}

func (s *Service) validateParticipantNotExists(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	_, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errors.Is(errors.NotFound, err) {
		return nil
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetServerParticipant")
	}
	return errors.E().
		Kind(errors.Exist).
		Messagef("user id='%d' already a participant of server id='%d'", userID.I64(), serverID.I64())
}

func (s *Service) validateTopicNameUnique(ctx context.Context, topic *domain.Topic) error {
	server, err := s.db.GetServer(ctx, topic.ServerID)
	if err != nil {
		return errors.E(err).Debug("s.db.GetServer")
	}

	for _, existingTopic := range server.Topics {
		if existingTopic.ID == topic.ID {
			continue
		}

		if existingTopic.Name == topic.Name && existingTopic.Type == topic.Type {
			return errors.E().
				Kind(errors.InvalidRequest).
				Messagef("topic %q with type %q already exists", topic.Name, topic.Type.String())
		}
	}

	return nil
}

func (s *Service) validateUserIsOwner(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	participant, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errors.Is(errors.NotFound, err) {
		return errors.E(err).
			Kind(errors.Permission).
			Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetServerParticipant")
	}

	if participant.Role != domain.ServerRoleOwner {
		return errors.E().
			Kind(errors.Permission).
			Message("only owners can perform this action")
	}

	return nil
}

func (s *Service) validateUserIsNotOwner(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	participant, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errors.Is(errors.NotFound, err) {
		return errors.E(err).Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errors.E(err).Debug("s.db.GetServerParticipant")
	}

	if participant.Role == domain.ServerRoleOwner {
		return errors.E().
			Kind(errors.Permission).
			Message("owners can not perform this action")
	}

	return nil
}

func (s *Service) validateTopicIsText(topic *domain.Topic) error {
	if topic.Type != domain.TopicTypeText {
		return errors.E().
			Kind(errors.InvalidRequest).
			Messagef("the topic id='%d' type must be text", topic.ID.I64())
	}
	return nil
}
