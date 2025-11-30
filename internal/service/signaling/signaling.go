package signaling

import (
	"chattery/internal/pb/api/websocketpb"
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type (
	MessageCallback func(ctx context.Context, msg *websocketpb.Message)
	CloseCallback   func(ctx context.Context)
)

type Service struct {
	conn          *websocket.Conn
	handlers      map[websocketpb.Type][]MessageCallback
	closeHandlers []CloseCallback
	mutex         sync.Mutex
}

func NewService(w http.ResponseWriter, r *http.Request) (*Service, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("upgrader.Upgrade: %w", err)
	}

	return &Service{
		conn:     conn,
		handlers: make(map[websocketpb.Type][]MessageCallback),
	}, nil
}

func (s *Service) RegisterMessageCallback(_type websocketpb.Type, callback MessageCallback) {
	s.handlers[_type] = append(s.handlers[_type], callback)
}

func (s *Service) RegisterCloseCallback(callback CloseCallback) {
	s.closeHandlers = append(s.closeHandlers, callback)
}

func (s *Service) Send(msg *websocketpb.Message) {
	msg.Id = rand.Text()

	data, err := protojson.Marshal(msg)
	if err != nil {
		slog.Error("[WebSocketManager] protojson.Marshal", slog.String("error", err.Error()))
		return
	}

	s.mutex.Lock()
	err = s.conn.WriteMessage(websocket.TextMessage, data)
	s.mutex.Unlock()

	if err != nil {
		slog.Error("[WebSocketManager] m.conn.WriteMessage", slog.String("error", err.Error()))
		return
	}
}

func (s *Service) ListenAndServe(ctx context.Context) {
	defer s.close(ctx)

	slog.Info("[WebSocketManager] start listening")

	for {
		select {
		case <-ctx.Done():
			slog.Warn("[WebSocketManager] context done", slog.String("error", ctx.Err().Error()))
			return
		default:
			if !s.readMessage(ctx) {
				return
			}
		}
	}
}

func (s *Service) readMessage(ctx context.Context) bool {
	s.mutex.Lock()
	mt, raw, err := s.conn.ReadMessage()
	s.mutex.Unlock()

	if err != nil {
		slog.Error("[WebSocketManager] m.conn.ReadMessage", slog.String("error", err.Error()))
		return false
	}

	if mt == websocket.CloseMessage {
		slog.Info("[WebSocketManager] closing connection")
		return false
	}

	fmt.Println(string(raw))

	message := new(websocketpb.Message)
	if err := protojson.Unmarshal(raw, message); err != nil {
		slog.Error("[WebSocketManager] protojson.Unmarshal", slog.String("error", err.Error()))
		return false
	}

	for _, callback := range s.handlers[message.GetType()] {
		callback(ctx, message)
	}

	return true
}

func (s *Service) close(ctx context.Context) {
	if err := s.conn.Close(); err != nil {
		slog.Error("[WebSocketManager] close", slog.String("error", err.Error()))
	}

	for _, callback := range s.closeHandlers {
		callback(ctx)
	}
}
