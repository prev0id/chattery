package server

import (
	"context"

	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

func (s *Service) validateServerExists(ctx context.Context, serverID domain.ServerID) error {
	_, err := s.db.GetServer(ctx, serverID)
	if errutil.Is(errutil.NotFound, err) {
		return errutil.E(err).Messagef("server id='%d' not found", serverID.I64())
	}

	if err != nil {
		return errutil.E(err).Debug("v.db.GetServer")
	}

	return nil
}

func (s *Service) validateTopicExists(ctx context.Context, topicID domain.TopicID) error {
	_, err := s.db.GetTopic(ctx, topicID)
	return err
}

func (s *Service) validateParticipantExists(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	_, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errutil.Is(errutil.NotFound, err) {
		return errutil.E(err).Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errutil.E(err).Debug("s.db.GetServerParticipant")
	}
	return nil
}

func (s *Service) validateParticipantNotExists(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	_, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errutil.Is(errutil.NotFound, err) {
		return nil
	}
	if err != nil {
		return errutil.E(err).Debug("s.db.GetServerParticipant")
	}
	return errutil.E().
		Kind(errutil.Exist).
		Messagef("user id='%d' already a participant of server id='%d'", userID.I64(), serverID.I64())
}

func (s *Service) validateTopicNameUnique(ctx context.Context, topic *domain.Topic) error {
	server, err := s.db.GetServer(ctx, topic.ServerID)
	if err != nil {
		return errutil.E(err).Debug("s.db.GetServer")
	}

	for _, existingTopic := range server.Topics {
		if existingTopic.ID == topic.ID {
			continue
		}

		if existingTopic.Name == topic.Name && existingTopic.Type == topic.Type {
			return errutil.E().
				Kind(errutil.InvalidRequest).
				Messagef("topic %q with type %q already exists", topic.Name, topic.Type.String())
		}
	}

	return nil
}

func (s *Service) validateUserIsOwner(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	participant, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errutil.Is(errutil.NotFound, err) {
		return errutil.E(err).
			Kind(errutil.Permission).
			Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errutil.E(err).Debug("s.db.GetServerParticipant")
	}

	if participant.Role != domain.ServerRoleOwner {
		return errutil.E().
			Kind(errutil.Permission).
			Message("only owners can perform this action")
	}

	return nil
}

func (s *Service) validateUserIsNotOwner(ctx context.Context, serverID domain.ServerID, userID domain.UserID) error {
	participant, err := s.db.GetServerParticipant(ctx, serverID, userID)
	if errutil.Is(errutil.NotFound, err) {
		return errutil.E(err).Messagef("user id='%d' not a participant of the server id='%d'", userID.I64(), serverID.I64())
	}
	if err != nil {
		return errutil.E(err).Debug("s.db.GetServerParticipant")
	}

	if participant.Role == domain.ServerRoleOwner {
		return errutil.E().
			Kind(errutil.Permission).
			Message("owners can not perform this action")
	}

	return nil
}
