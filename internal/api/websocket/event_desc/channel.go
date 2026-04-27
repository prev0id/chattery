package event_desc

type Channel struct {
	Type ChannelType `json:"type"`
	ID   int64       `json:"id"`
}

type ChannelType string

const (
	ChannelDM         ChannelType = "dm"
	ChannelTextTopic  ChannelType = "text_topic"
	ChannelVoiceTopic ChannelType = "voice_topic"
)
