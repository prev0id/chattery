package server_adapter

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

func convertServersFromDB(rows []*postgres.GetUserServersRow) []*domain.Server {
	if len(rows) == 0 {
		return nil
	}

	grouped := sliceutil.GroupBy(rows, func(row *postgres.GetUserServersRow) int64 {
		return row.ID
	})

	servers := make([]*domain.Server, 0, len(grouped))
	for serverID, serverRows := range grouped {
		if len(serverRows) == 0 {
			continue
		}

		server := &domain.Server{
			ID:   domain.ServerID(serverID),
			Name: serverRows[0].Name,
		}

		for _, row := range serverRows {
			if row.TopicID.Valid {
				topic := &domain.Topic{
					ID:       domain.TopicID(row.TopicID.Int64),
					ServerID: domain.ServerID(serverID),
					Name:     row.TopicName.String,
					Type:     domain.TopicType(row.TopicType.String),
				}
				server.Topics = append(server.Topics, *topic)
			}
		}

		servers = append(servers, server)
	}

	return servers
}
