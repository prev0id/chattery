package domain

import "time"

type Server struct {
	ID       ServerID
	Name     string
	JoinedAt time.Time
	Topics   []*Topic
}

type ServerID int64

func (id ServerID) I64() int64 { return int64(id) }

type Topic struct {
	ID        TopicID
	ServerID  ServerID
	Name      string
	Type      TopicType
	CreatedAt time.Time
}

type TopicID int64

func (id TopicID) I64() int64 { return int64(id) }

type TopicType string

func (t TopicType) String() string { return string(t) }

const (
	TopicTypeText  TopicType = "text"
	TopicTypeVoice TopicType = "voice"
)

type ServerParticipant struct {
	UserID   UserID
	ServerID ServerID
	Role     ServerRole
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
	ID        TopicMessageID
	TopicID   TopicID
	SenderID  UserID
	Text      string
	CreatedAt time.Time
}
