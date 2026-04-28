package server

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"
	"chattery/internal/utils/sliceutil"
)

func convertServerMessageFromDB(msg *postgres.TopicMessage) *domain.TopicMessage {
	return &domain.TopicMessage{
		ID:        domain.TopicMessageID(msg.ID),
		TopicID:   domain.TopicID(msg.TopicID),
		SenderID:  domain.UserID(msg.UserID),
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

func convertServerWithTopicsFromDB(rows []*postgres.GetServerRow) *domain.Server {
	if len(rows) == 0 {
		return nil
	}

	server := &domain.Server{
		ID:   domain.ServerID(rows[0].ID),
		Name: rows[0].Name,
	}

	for _, row := range rows {
		if row.TopicID.Valid {
			topic := &domain.Topic{
				ID:        domain.TopicID(row.TopicID.Int64),
				ServerID:  domain.ServerID(row.ID),
				Name:      row.TopicName.String,
				Type:      domain.TopicType(row.TopicType.String),
				CreatedAt: row.TopicCreatedAt.Time,
			}
			server.Topics = append(server.Topics, topic)
		}
	}

	return server
}

func convertTopicFromDB(topic *postgres.Topic) *domain.Topic {
	return &domain.Topic{
		ID:       domain.TopicID(topic.ID),
		ServerID: domain.ServerID(topic.ServerID),
		Name:     topic.Name,
		Type:     domain.TopicType(topic.Type),
	}
}

func convertServerParticipantFromDB(participant *postgres.ServerParticipant) *domain.ServerParticipant {
	return &domain.ServerParticipant{
		UserID:   domain.UserID(participant.UserID),
		ServerID: domain.ServerID(participant.ServerID),
		Role:     domain.ServerRole(participant.Role),
	}
}

func convertServersFromDB(rows []*postgres.GetUserServersRow) []*domain.Server {
	if len(rows) == 0 {
		return nil
	}

	grouped := sliceutil.GroupBy(rows, func(row *postgres.GetUserServersRow) int64 {
		return row.ID
	})

	servers := make([]*domain.Server, 0, len(grouped))
	for serverID, serverRows := range grouped {
		if server := convertUserServerFromRows(serverID, serverRows); server != nil {
			servers = append(servers, server)
		}
	}

	return servers
}

func convertUserServerFromRows(serverID int64, serverRows []*postgres.GetUserServersRow) *domain.Server {
	if len(serverRows) == 0 {
		return nil
	}

	server := &domain.Server{
		ID:       domain.ServerID(serverID),
		Name:     serverRows[0].Name,
		JoinedAt: serverRows[0].JoinedAt,
		Role:     domain.ServerRole(serverRows[0].Role),
	}

	for _, row := range serverRows {
		if topic := convertTopicFromUserServerRow(serverID, row); topic != nil {
			server.Topics = append(server.Topics, topic)
		}
	}

	return server
}

func convertServersFromListDB(rows []*postgres.ListServersRow) []*domain.Server {
	if len(rows) == 0 {
		return nil
	}

	grouped := sliceutil.GroupBy(rows, func(row *postgres.ListServersRow) int64 {
		return row.ID
	})

	servers := make([]*domain.Server, 0, len(grouped))
	for serverID, serverRows := range grouped {
		if server := convertServerFromListRows(serverID, serverRows); server != nil {
			servers = append(servers, server)
		}
	}

	return servers
}

func convertServerFromListRows(serverID int64, serverRows []*postgres.ListServersRow) *domain.Server {
	if len(serverRows) == 0 {
		return nil
	}

	server := &domain.Server{
		ID:   domain.ServerID(serverID),
		Name: serverRows[0].Name,
	}

	for _, row := range serverRows {
		if topic := convertTopicFromRow(serverID, row); topic != nil {
			server.Topics = append(server.Topics, topic)
		}
	}

	return server
}

func convertTopicFromRow(serverID int64, row *postgres.ListServersRow) *domain.Topic {
	if !row.TopicID.Valid {
		return nil
	}
	return &domain.Topic{
		ID:        domain.TopicID(row.TopicID.Int64),
		ServerID:  domain.ServerID(serverID),
		Name:      row.TopicName.String,
		Type:      domain.TopicType(row.TopicType.String),
		CreatedAt: row.TopicCreatedAt.Time,
	}
}

func convertTopicFromUserServerRow(serverID int64, row *postgres.GetUserServersRow) *domain.Topic {
	if !row.TopicID.Valid {
		return nil
	}
	return &domain.Topic{
		ID:        domain.TopicID(row.TopicID.Int64),
		ServerID:  domain.ServerID(serverID),
		Name:      row.TopicName.String,
		Type:      domain.TopicType(row.TopicType.String),
		CreatedAt: row.TopicCreatedAt.Time,
	}
}
