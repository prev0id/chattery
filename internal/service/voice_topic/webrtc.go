package voice_topic

import (
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"chattery/internal/config"
)

type webrtcAPI struct {
	api           *webrtc.API
	configuration webrtc.Configuration
}

func newWebRTCAPI(cfg *config.Config) *webrtcAPI {
	mediaEngine := &webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		panic(err)
	}

	interceptorRegistry := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
		panic(err)
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
	}
}

func (w *webrtcAPI) newPeerConnection() (*webrtc.PeerConnection, error) {
	return w.api.NewPeerConnection(w.configuration)
}
