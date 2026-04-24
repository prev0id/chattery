package domain

import (
	"strconv"
	"time"
)

type Server struct {
	JoinedAt time.Time
	Name     string
	Role     ServerRole
	Topics   []*Topic
	ID       ServerID
}

type ServerID int64

func (id ServerID) I64() int64 { return int64(id) }

type Topic struct {
	CreatedAt time.Time
	Name      string
	Type      TopicType
	ID        TopicID
	ServerID  ServerID
}

type TopicID int64

func (id TopicID) I64() int64 { return int64(id) }

func (id TopicID) String() string { return strconv.FormatInt(id.I64(), 10) }

type TopicType string

func (t TopicType) String() string { return string(t) }

const (
	TopicTypeText  TopicType = "text"
	TopicTypeVoice TopicType = "voice"
)

type ServerParticipant struct {
	Role     ServerRole
	UserID   UserID
	ServerID ServerID
}

type ServerRole string

const (
	ServerRoleOwner  ServerRole = "owner"
	ServerRoleMember ServerRole = "member"
)

func (role ServerRole) String() string { return string(role) }

type TopicMessageID int64

func (id TopicMessageID) I64() int64 { return int64(id) }

type TopicMessage struct {
	CreatedAt time.Time
	Text      string
	ID        TopicMessageID
	TopicID   TopicID
	SenderID  UserID
}
