package voice_topic

import (
	"time"

	"chattery/internal/domain"
)

func (s *Service) ICEServers(userID domain.UserID) ([]domain.VoiceICEServer, error) {
	return s.webrtc.iceServers(userID, time.Now()), nil
}
