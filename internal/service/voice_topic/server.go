package voice_topic

import (
	"encoding/json"
	"sync"

	"github.com/pion/webrtc/v4"

	"chattery/internal/domain"
	"chattery/internal/utils/errors"
)

var (
	codecs = []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio}

	directionReceiveOnly = webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	}
)

type WsEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Topic struct {
	peers       map[domain.UserID]*webrtc.PeerConnection
	tracks      map[string]*webrtc.TrackLocalStaticRTP
	peersMutex  sync.RWMutex
	tracksMutex sync.RWMutex
}

func (t *Topic) addPeerConnection(conn *webrtc.PeerConnection, userID domain.UserID) {
	t.peersMutex.Lock()
	defer t.peersMutex.Unlock()

	t.peers[userID] = conn
}

func (t *Topic) onTrack(userID domain.UserID) trackCallback {
	return func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		t.peersMutex.RLock()
		defer t.peersMutex.RUnlock()
	}
}

func (t *Topic) addTrack(track *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
	local, err := webrtc.NewTrackLocalStaticRTP(
		track.Codec().RTPCodecCapability,
		track.ID(),
		track.StreamID(),
	)
	if err != nil {
		return nil, errors.E(err).
			Debug("webrtc.NewTrackLocalStaticRTP", track.ID(), track.StreamID()).
			Message("Unable to create track local RTP, try again later")
	}

	t.tracksMutex.Lock()
	defer t.tracksMutex.Unlock()
	// TODO add signaling

	t.tracks[track.ID()] = local

	return local, nil
}

func (t *Topic) removeTrack(local *webrtc.TrackRemote) {
	t.tracksMutex.Lock()
	defer t.tracksMutex.Unlock()
	// TODO add signaling

	delete(t.tracks, local.ID())
}

func (t *Topic) signalPeers() {

}

type Service struct {
	topics map[domain.TopicID]*Topic
	m      sync.RWMutex
}

func (s *Service) Register(topicID domain.TopicID, userID domain.UserID) error {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return errors.
			E(err).
			Debug("webrtc.NewPeerConnection").
			Message("Unable to create peer connection, try again later")
	}

	for _, codec := range codecs {
		if _, err := peerConnection.AddTransceiverFromKind(codec, directionReceiveOnly); err != nil {
			return errors.
				E(err).
				Debug("peerConnection.AddTransceiverFromKind").
				Message("Unable to create RTC transceiver, try again later")
		}
	}

	topic := s.getTopic(topicID)

	peerConnection.OnICECandidate(s.onICECandidate())
	peerConnection.OnConnectionStateChange(s.onConnectionStateChange())
	peerConnection.OnTrack(topic.onTrack(userID))
	peerConnection.OnICEConnectionStateChange(s.onICEConnectionStateChange())

	topic.addPeerConnection(peerConnection, userID)

	return nil
}

func (s *Service) getTopic(topicID domain.TopicID) *Topic {
	s.m.Lock()
	defer s.m.Unlock()

	topic, exists := s.topics[topicID]
	if !exists {
		topic := &Topic{peers: make(map[domain.UserID]*webrtc.PeerConnection)}
		s.topics[topicID] = topic
	}
	return topic
}
