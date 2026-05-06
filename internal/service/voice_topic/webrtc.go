package voice_topic

import (
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"chattery/internal/config"
	"chattery/internal/utils/errutil"
)

const maxUDPPort = 65535

type webrtcAPI struct {
	api           *webrtc.API
	configuration webrtc.Configuration
}

func newWebRTCAPI(cfg *config.Config) (*webrtcAPI, error) {
	mediaEngine := &webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		return nil, errutil.E(err).Debug("mediaEngine.RegisterDefaultCodecs")
	}

	interceptorRegistry := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
		return nil, errutil.E(err).Debug("webrtc.RegisterDefaultInterceptors")
	}

	configuration := webrtc.Configuration{}
	if cfg.Voice.STUNServer != "" {
		configuration.ICEServers = []webrtc.ICEServer{
			{URLs: []string{cfg.Voice.STUNServer}},
		}
	}

	settingEngine, err := newSettingEngine(cfg)
	if err != nil {
		return nil, err
	}

	return &webrtcAPI{
		api: webrtc.NewAPI(
			webrtc.WithMediaEngine(mediaEngine),
			webrtc.WithInterceptorRegistry(interceptorRegistry),
			webrtc.WithSettingEngine(settingEngine),
		),
		configuration: configuration,
	}, nil
}

func newSettingEngine(cfg *config.Config) (webrtc.SettingEngine, error) {
	settingEngine := webrtc.SettingEngine{}
	if err := setUDPPortRange(&settingEngine, cfg); err != nil {
		return settingEngine, err
	}
	if err := setPublicIP(&settingEngine, cfg); err != nil {
		return settingEngine, err
	}
	return settingEngine, nil
}

func setUDPPortRange(settingEngine *webrtc.SettingEngine, cfg *config.Config) error {
	if cfg.Voice.ICEUDPPortMin == 0 && cfg.Voice.ICEUDPPortMax == 0 {
		return nil
	}
	if cfg.Voice.ICEUDPPortMin <= 0 || cfg.Voice.ICEUDPPortMax <= 0 {
		return errutil.E().Message("voice ICE UDP port range must set both min and max")
	}
	if cfg.Voice.ICEUDPPortMin > maxUDPPort || cfg.Voice.ICEUDPPortMax > maxUDPPort {
		return errutil.E().Messagef("voice ICE UDP port range must be <= %d", maxUDPPort)
	}
	if err := settingEngine.SetEphemeralUDPPortRange(
		uint16(cfg.Voice.ICEUDPPortMin),
		uint16(cfg.Voice.ICEUDPPortMax),
	); err != nil {
		return errutil.E(err).Debug("settingEngine.SetEphemeralUDPPortRange")
	}
	return nil
}

func setPublicIP(settingEngine *webrtc.SettingEngine, cfg *config.Config) error {
	if cfg.Voice.ICEPublicIP == "" {
		return nil
	}
	if err := settingEngine.SetICEAddressRewriteRules(webrtc.ICEAddressRewriteRule{
		External:        []string{cfg.Voice.ICEPublicIP},
		AsCandidateType: webrtc.ICECandidateTypeHost,
		Mode:            webrtc.ICEAddressRewriteReplace,
	}); err != nil {
		return errutil.E(err).Debug("settingEngine.SetICEAddressRewriteRules")
	}
	return nil
}

func (w *webrtcAPI) newPeerConnection() (*webrtc.PeerConnection, error) {
	return w.api.NewPeerConnection(w.configuration)
}
