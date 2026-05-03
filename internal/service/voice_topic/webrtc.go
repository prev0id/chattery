package voice_topic

import (
	"fmt"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"chattery/internal/config"
)

type webrtcAPI struct {
	api           *webrtc.API
	configuration webrtc.Configuration
}

func newWebRTCAPI(cfg *config.Config) (*webrtcAPI, error) {
	mediaEngine := &webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		return nil, fmt.Errorf("mediaEngine.RegisterDefaultCodecs: %w", err)
	}

	interceptorRegistry := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
		return nil, fmt.Errorf("webrtc.RegisterDefaultInterceptors: %w", err)
	}

	configuration := webrtc.Configuration{}
	if cfg.Voice.STUNServer != "" {
		configuration.ICEServers = []webrtc.ICEServer{
			{URLs: []string{cfg.Voice.STUNServer}},
		}
	}

	return &webrtcAPI{
		api: webrtc.NewAPI(
			webrtc.WithMediaEngine(mediaEngine),
			webrtc.WithInterceptorRegistry(interceptorRegistry),
		),
		configuration: configuration,
	}, nil
}

func (w *webrtcAPI) newPeerConnection() (*webrtc.PeerConnection, error) {
	return w.api.NewPeerConnection(w.configuration)
}
