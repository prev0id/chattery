package voice_topic

import (
	"github.com/pion/webrtc/v4"
)

type (
	iceCandidateCallback          = func(candidate *webrtc.ICECandidate)
	connectionStateChangeCallback = func(state webrtc.PeerConnectionState)
	trackCallback                 = func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)
	iceConnectionStateChange      = func(state webrtc.ICEConnectionState)
)

func (s *Service) onICECandidate() iceCandidateCallback {
	return func(candidate *webrtc.ICECandidate) {}
}

func (s *Service) onConnectionStateChange() connectionStateChangeCallback {
	return func(state webrtc.PeerConnectionState) {}
}

func (s *Service) onICEConnectionStateChange() iceConnectionStateChange {
	return func(state webrtc.ICEConnectionState) {}
}
