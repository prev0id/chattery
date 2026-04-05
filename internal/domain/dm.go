package domain

import "time"

type DMID int64

func (id DMID) I64() int64 { return int64(id) }

type DM struct {
	ID                 DMID
	LastActivityAt     time.Time
	LastMessage        DMMessage
	OtherParticipantID UserID
}

type DMParticipant struct {
	UserID            UserID
	DMID              DMID
	LastReadMessageID DMMessageID
}

type DMMessageID int64

func (id DMMessageID) I64() int64 { return int64(id) }

type DMMessage struct {
	ID        DMMessageID
	DMID      DMID
	SenderID  UserID
	Text      string
	CreatedAt time.Time
}
