package voice_topic

import (
	"crypto/hmac"
	"crypto/sha1" // #nosec G505 -- coturn REST API requires HMAC-SHA1 credentials.
	"encoding/base64"
	"strconv"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v4"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errutil"
)

const (
	maxUDPPort               = 65535
	defaultTURNCredentialTTL = 10 * time.Minute
)

type webrtcAPI struct {
	api               *webrtc.API
	turnAuthSecret    string
	stunServers       []string
	turnServers       []string
	turnCredentialTTL time.Duration
}

func newWebRTCAPI(cfg *config.Config) (*webrtcAPI, error) {
	if len(cfg.Voice.TURNServers) > 0 && cfg.Voice.TURNAuthSecret == "" {
		return nil, errutil.E().
			Message("voice TURN auth secret is required when TURN servers are configured")
	}

	mediaEngine := &webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		return nil, errutil.E(err).Debug("mediaEngine.RegisterDefaultCodecs")
	}

	interceptorRegistry := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(mediaEngine, interceptorRegistry); err != nil {
		return nil, errutil.E(err).Debug("webrtc.RegisterDefaultInterceptors")
	}

	settingEngine, err := newSettingEngine(cfg)
	if err != nil {
		return nil, err
	}
	turnCredentialTTL := cfg.Voice.TURNCredentialTTL
	if turnCredentialTTL <= 0 {
		turnCredentialTTL = defaultTURNCredentialTTL
	}

	return &webrtcAPI{
		api: webrtc.NewAPI(
			webrtc.WithMediaEngine(mediaEngine),
			webrtc.WithInterceptorRegistry(interceptorRegistry),
			webrtc.WithSettingEngine(settingEngine),
		),
		stunServers:       cfg.Voice.STUNServers,
		turnServers:       cfg.Voice.TURNServers,
		turnAuthSecret:    cfg.Voice.TURNAuthSecret,
		turnCredentialTTL: turnCredentialTTL,
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

func (w *webrtcAPI) newPeerConnection(userID domain.UserID) (*webrtc.PeerConnection, error) {
	configuration := webrtc.Configuration{
		ICEServers: w.pionICEServers(userID),
	}
	return w.api.NewPeerConnection(configuration)
}

func (w *webrtcAPI) pionICEServers(userID domain.UserID) []webrtc.ICEServer {
	iceServers := w.iceServers(userID, time.Now())
	result := make([]webrtc.ICEServer, 0, len(iceServers))
	for _, iceServer := range iceServers {
		result = append(result, webrtc.ICEServer{
			URLs:       iceServer.URLs,
			Username:   iceServer.Username,
			Credential: iceServer.Credential,
		})
	}
	return result
}

func (w *webrtcAPI) iceServers(userID domain.UserID, now time.Time) []domain.VoiceICEServer {
	iceServers := make([]domain.VoiceICEServer, 0, 2)
	if len(w.stunServers) > 0 {
		iceServers = append(iceServers, domain.VoiceICEServer{
			URLs: w.stunServers,
		})
	}
	if len(w.turnServers) > 0 {
		username, credential := w.turnCredentials(userID, now)
		iceServers = append(iceServers, domain.VoiceICEServer{
			URLs:       w.turnServers,
			Username:   username,
			Credential: credential,
		})
	}
	return iceServers
}

func (w *webrtcAPI) turnCredentials(userID domain.UserID, now time.Time) (username string, credential string) {
	expiresAtSec := now.Add(w.turnCredentialTTL).Unix()
	username = strconv.FormatInt(expiresAtSec, 10) + ":" + userID.String()

	mac := hmac.New(sha1.New, []byte(w.turnAuthSecret))
	_, _ = mac.Write([]byte(username))
	credential = base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return username, credential
}
