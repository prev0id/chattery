package domain

import (
	"context"

	"chattery/internal/api/websocket/event_desc"
)

type Connection interface {
	ReadPump(ctx context.Context)
	WritePump(ctx context.Context)
	WriteEvent(ctx context.Context, event *event_desc.Event)
}
