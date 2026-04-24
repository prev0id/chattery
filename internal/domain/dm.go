package domain

import "time"

type DMID int64

func (id DMID) I64() int64 { return int64(id) }

type DM struct {
	LastActivityAt     time.Time
	LastMessage        DMMessage
	ID                 DMID
	OtherParticipantID UserID
	LastReadMessageID  DMMessageID
}

type DMParticipant struct {
	UserID            UserID
	DMID              DMID
	LastReadMessageID DMMessageID
}

type DMMessageID int64

func (id DMMessageID) I64() int64 { return int64(id) }

type DMMessage struct {
	CreatedAt time.Time
	Text      string
	ID        DMMessageID
	DMID      DMID
	SenderID  UserID
}
