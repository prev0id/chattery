package domain

type UserMessageType string

const (
	UserMessageTypeDM    UserMessageType = "dm"
	UserMessageTypeTopic UserMessageType = "topic"
)

type UserMessage struct {
	DMMessage *DMMessage      `json:"dm_message,omitempty"`
	TopicMsg  *TopicMessage   `json:"topic_message,omitempty"`
	Type      UserMessageType `json:"type"`
	ChannelID int64           `json:"channel_id"`
}
