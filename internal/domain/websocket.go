package domain

import "context"

type ChannelType string

const (
	ChannelDM        ChannelType = "dm"
	ChannelTextTopic ChannelType = "text_topic"
	ChannelVoice     ChannelType = "voice"
)

type Connection interface {
	ReadPump(ctx context.Context)
	WritePump(ctx context.Context)
}
