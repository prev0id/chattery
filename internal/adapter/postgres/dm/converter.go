package dm_adapter

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"
)

func convertDMMessageFromDB(msg *postgres.DmMessage) *domain.DMMessage {
	return &domain.DMMessage{
		ID:        domain.DMMessageID(msg.ID),
		DMID:      domain.DMID(msg.DmID),
		SenderID:  domain.UserID(msg.UserID),
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

func convertDMFromDB(dm *postgres.UserDMsRow) *domain.DM {
	var lastMessage domain.DMMessage
	if dm.LastMessageID.Valid {
		lastMessage = domain.DMMessage{
			ID:        domain.DMMessageID(dm.LastMessageID.Int64),
			DMID:      domain.DMID(dm.DmID),
			SenderID:  domain.UserID(dm.LastMessageSenderID.Int64),
			Text:      dm.LastMessageText.String,
			CreatedAt: dm.LastMessageCreatedAt.Time,
		}
	}

	return &domain.DM{
		ID:                 domain.DMID(dm.DmID),
		LastActivityAt:     dm.LastActivityAt,
		LastMessage:        lastMessage,
		OtherParticipantID: domain.UserID(dm.OtherParticipantID),
	}
}

func convertDMParticipantFromDB(participant *postgres.DmParticipant) *domain.DMParticipant {
	return &domain.DMParticipant{
		DMID:              domain.DMID(participant.DmID),
		UserID:            domain.UserID(participant.UserID),
		LastReadMessageID: domain.DMMessageID(participant.LastReadMessageID),
	}
}
