package dm

import (
	"chattery/internal/client/postgres"
	"chattery/internal/domain"

	"github.com/jackc/pgx/v5/pgtype"
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
	return &domain.DM{
		ID:                 domain.DMID(dm.DmID),
		OtherParticipantID: domain.UserID(dm.OtherParticipantID),
		LastReadMessageID:  domain.DMMessageID(dm.LastReadMessageID),
		LastActivityAt:     dm.LastActivityAt,
		LastMessage:        convertLastDMMessageFromDB(dm.DmID, dm.LastMessageID, dm.LastMessageSenderID, dm.LastMessageText, dm.LastMessageCreatedAt),
	}
}

func convertDMFromListDB(dm *postgres.ListDMsRow) *domain.DM {
	return &domain.DM{
		ID:             domain.DMID(dm.DmID),
		LastActivityAt: dm.LastActivityAt,
		LastMessage:    convertLastDMMessageFromDB(dm.DmID, dm.LastMessageID, dm.LastMessageSenderID, dm.LastMessageText, dm.LastMessageCreatedAt),
	}
}

func convertLastDMMessageFromDB(
	dmID int64,
	messageID pgtype.Int8,
	senderID pgtype.Int8,
	text pgtype.Text,
	createdAt pgtype.Timestamp,
) domain.DMMessage {
	if !messageID.Valid {
		return domain.DMMessage{}
	}

	return domain.DMMessage{
		ID:        domain.DMMessageID(messageID.Int64),
		DMID:      domain.DMID(dmID),
		SenderID:  domain.UserID(senderID.Int64),
		Text:      text.String,
		CreatedAt: createdAt.Time,
	}
}

func convertDMParticipantFromDB(participant *postgres.DmParticipant) *domain.DMParticipant {
	return &domain.DMParticipant{
		DMID:              domain.DMID(participant.DmID),
		UserID:            domain.UserID(participant.UserID),
		LastReadMessageID: domain.DMMessageID(participant.LastReadMessageID),
	}
}
