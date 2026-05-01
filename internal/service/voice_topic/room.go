package voice_topic

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"

	"chattery/internal/api/websocket/event_desc"
	"chattery/internal/domain"
	"chattery/internal/utils/bind"
	"chattery/internal/utils/errutil"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/render"
)

const iceBatchDelay = 50 * time.Millisecond

type topic struct {
	participants map[domain.UserID]*participant
	tracks       map[string]*localTrack
	service      *Service
	id           domain.TopicID
	m            sync.RWMutex
}

type participant struct {
	senders       map[string]*webrtc.RTPSender
	pc            *webrtc.PeerConnection
	iceBatchTimer *time.Timer
	pendingICE    []webrtc.ICECandidateInit
	iceBatch      []event_desc.VoiceICECandidatePayload
	m             sync.Mutex
	userID        domain.UserID
	remoteReady   bool
	renegotiating bool
}

type localTrack struct {
	publisher *participant
	track     *webrtc.TrackLocalStaticRTP
	cancel    context.CancelFunc
	id        string
	owner     domain.UserID
	ssrc      uint32
	kind      webrtc.RTPCodecType
}

func newTopic(topicID domain.TopicID, service *Service) *topic {
	return &topic{
		service:      service,
		id:           topicID,
		participants: make(map[domain.UserID]*participant),
		tracks:       make(map[string]*localTrack),
	}
}

func (r *topic) addParticipant(userID domain.UserID) bool {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.participants[userID]; ok {
		return false
	}
	r.participants[userID] = &participant{
		userID:  userID,
		senders: make(map[string]*webrtc.RTPSender),
	}
	return true
}

func (r *topic) removeParticipant(userID domain.UserID) bool {
	r.m.Lock()
	participant, ok := r.participants[userID]
	if !ok {
		r.m.Unlock()
		return false
	}
	delete(r.participants, userID)

	ownedTracks := make([]*localTrack, 0)
	for key, track := range r.tracks {
		if track.owner == userID {
			ownedTracks = append(ownedTracks, track)
			delete(r.tracks, key)
		}
	}
	r.m.Unlock()

	for _, track := range ownedTracks {
		track.cancel()
		r.removeTrackFromParticipants(track.id, userID, true)
	}
	if participant.pc != nil {
		if err := participant.pc.Close(); err != nil {
			logger.Error(err, "[voice_topic] participant pc.Close")
		}
	}
	return true
}

func (r *topic) hasParticipant(userID domain.UserID) bool {
	r.m.RLock()
	defer r.m.RUnlock()

	_, ok := r.participants[userID]
	return ok
}

func (r *topic) empty() bool {
	r.m.RLock()
	defer r.m.RUnlock()

	return len(r.participants) == 0
}

func (r *topic) participantIDs() []domain.UserID {
	r.m.RLock()
	defer r.m.RUnlock()

	participants := make([]domain.UserID, 0, len(r.participants))
	for participant := range r.participants {
		participants = append(participants, participant)
	}
	return participants
}

func (r *topic) handleOffer(ctx context.Context, userID domain.UserID, payload json.RawMessage) error {
	desc, bindErr := bind.JSONBytes[event_desc.VoiceSessionDescriptionPayload](payload)
	if bindErr != nil {
		return bindErr
	}

	participant, peerErr := r.ensurePeerConnection(ctx, userID)
	if peerErr != nil {
		return peerErr
	}

	if setRemoteErr := participant.pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  desc.SDP,
	}); setRemoteErr != nil {
		return errutil.E(setRemoteErr).Debug("participant.pc.SetRemoteDescription")
	}

	participant.m.Lock()
	participant.remoteReady = true
	pendingICE := participant.pendingICE
	participant.pendingICE = nil
	participant.m.Unlock()

	for _, candidate := range pendingICE {
		if addCandidateErr := participant.pc.AddICECandidate(candidate); addCandidateErr != nil {
			logger.ErrorCtx(ctx, errutil.E(addCandidateErr).Debug("participant.pc.AddICECandidate"), "[voice_topic] AddICECandidate")
		}
	}

	if addTrackErr := r.addExistingTracks(participant); addTrackErr != nil {
		return addTrackErr
	}

	answer, answerErr := participant.pc.CreateAnswer(nil)
	if answerErr != nil {
		return errutil.E(answerErr).Debug("participant.pc.CreateAnswer")
	}

	if setLocalErr := participant.pc.SetLocalDescription(answer); setLocalErr != nil {
		return errutil.E(setLocalErr).Debug("participant.pc.SetLocalDescription")
	}

	if err := r.sendSessionDescription(ctx, userID, event_desc.TypeVoiceAnswer, answer); err != nil {
		return err
	}
	r.requestParticipantKeyFrames(participant)
	return nil
}

func (r *topic) handleAnswer(_ context.Context, userID domain.UserID, payload json.RawMessage) error {
	desc, err := bind.JSONBytes[event_desc.VoiceSessionDescriptionPayload](payload)
	if err != nil {
		return err
	}

	participant := r.getParticipant(userID)
	if participant == nil || participant.pc == nil {
		return errutil.E().Kind(errutil.InvalidRequest).Message("voice peer connection is not ready")
	}

	if err := participant.pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  desc.SDP,
	}); err != nil {
		return errutil.E(err).Debug("participant.pc.SetRemoteDescription")
	}

	participant.m.Lock()
	participant.remoteReady = true
	participant.renegotiating = false
	participant.m.Unlock()

	r.requestParticipantKeyFrames(participant)

	return nil
}

func (r *topic) removeOwnedTracksByKind(owner domain.UserID, kind webrtc.RTPCodecType, shouldRenegotiate bool) {
	r.m.Lock()
	ownedTracks := make([]*localTrack, 0)
	for key, track := range r.tracks {
		if track.owner == owner && track.kind == kind {
			ownedTracks = append(ownedTracks, track)
			delete(r.tracks, key)
		}
	}
	r.m.Unlock()

	for _, track := range ownedTracks {
		track.cancel()
		r.removeTrackFromParticipants(track.id, owner, shouldRenegotiate)
	}
}

func (r *topic) handleICECandidate(ctx context.Context, userID domain.UserID, payload json.RawMessage) error {
	candidate, err := bind.JSONBytes[event_desc.VoiceICECandidatePayload](payload)
	if err != nil {
		return err
	}

	participant, err := r.ensurePeerConnection(ctx, userID)
	if err != nil {
		return err
	}

	init := webrtc.ICECandidateInit{
		Candidate:        candidate.Candidate,
		SDPMid:           candidate.SDPMid,
		SDPMLineIndex:    candidate.SDPMLineIndex,
		UsernameFragment: candidate.UsernameFragment,
	}

	participant.m.Lock()
	remoteReady := participant.remoteReady
	if !remoteReady {
		participant.pendingICE = append(participant.pendingICE, init)
		participant.m.Unlock()
		return nil
	}
	participant.m.Unlock()

	if err := participant.pc.AddICECandidate(init); err != nil {
		return errutil.E(err).Debug("participant.pc.AddICECandidate")
	}
	return nil
}

func (r *topic) handleICECandidates(ctx context.Context, userID domain.UserID, payload json.RawMessage) error {
	candidates, err := bind.JSONBytes[event_desc.VoiceICECandidatesPayload](payload)
	if err != nil {
		return err
	}

	for _, candidate := range candidates.Candidates {
		if candidate.Candidate == "" {
			continue
		}

		raw, err := render.JSONBytes(candidate)
		if err != nil {
			return errutil.E(err).Debug("render.JSONBytes")
		}
		if err := r.handleICECandidate(ctx, userID, raw); err != nil {
			return err
		}
	}
	return nil
}

func (r *topic) ensurePeerConnection(ctx context.Context, userID domain.UserID) (*participant, error) {
	participant := r.getParticipant(userID)
	if participant == nil {
		r.addParticipant(userID)
		participant = r.getParticipant(userID)
	}
	if participant.pc != nil {
		return participant, nil
	}

	pc, err := r.service.webrtc.newPeerConnection()
	if err != nil {
		return nil, errutil.E(err).Debug("webrtc.NewPeerConnection")
	}

	participant.pc = pc
	pc.OnICECandidate(r.onICECandidate(userID))
	pc.OnTrack(r.onTrack(ctx, userID))
	pc.OnConnectionStateChange(r.onConnectionStateChange(ctx, userID))

	return participant, nil
}

func (r *topic) getParticipant(userID domain.UserID) *participant {
	r.m.RLock()
	defer r.m.RUnlock()
	return r.participants[userID]
}

func (r *topic) addExistingTracks(participant *participant) error {
	r.m.RLock()
	tracks := make([]*localTrack, 0, len(r.tracks))
	for _, track := range r.tracks {
		if track.owner != participant.userID {
			tracks = append(tracks, track)
		}
	}
	r.m.RUnlock()

	for _, track := range tracks {
		if err := addTrackToParticipant(participant, track); err != nil {
			return err
		}
	}
	return nil
}

func addTrackToParticipant(participant *participant, track *localTrack) error {
	participant.m.Lock()
	if participant.pc == nil || !participant.remoteReady {
		participant.m.Unlock()
		return nil
	}
	if _, ok := participant.senders[track.id]; ok {
		participant.m.Unlock()
		return nil
	}

	sender, err := participant.pc.AddTrack(track.track)
	if err != nil {
		participant.m.Unlock()
		return errutil.E(err).Debug("participant.pc.AddTrack")
	}
	participant.senders[track.id] = sender
	participant.m.Unlock()

	go readRTCP(sender, track)
	requestKeyFrameBurst(track)
	return nil
}

func (r *topic) removeTrackFromParticipants(trackID string, owner domain.UserID, shouldRenegotiate bool) {
	r.m.RLock()
	participants := make([]*participant, 0, len(r.participants))
	for _, participant := range r.participants {
		if participant.userID != owner {
			participants = append(participants, participant)
		}
	}
	r.m.RUnlock()

	for _, participant := range participants {
		removeTrackFromParticipant(participant, trackID)
		if shouldRenegotiate {
			r.renegotiate(context.Background(), participant)
		}
	}
}

func removeTrackFromParticipant(participant *participant, trackID string) {
	participant.m.Lock()
	sender := participant.senders[trackID]
	delete(participant.senders, trackID)
	participant.m.Unlock()

	if sender == nil || participant.pc == nil {
		return
	}
	if err := participant.pc.RemoveTrack(sender); err != nil {
		logger.Error(errutil.E(err).Debug("participant.pc.RemoveTrack"), "[voice_topic] RemoveTrack")
	}
}

func (r *topic) onTrack(ctx context.Context, userID domain.UserID) func(*webrtc.TrackRemote, *webrtc.RTPReceiver) {
	return func(remote *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		if remote.Kind() == webrtc.RTPCodecTypeVideo {
			r.removeOwnedTracksByKind(userID, webrtc.RTPCodecTypeVideo, false)
		}

		track, err := webrtc.NewTrackLocalStaticRTP(
			remote.Codec().RTPCodecCapability,
			remote.ID(),
			strconv.FormatInt(userID.I64(), 10),
		)
		if err != nil {
			logger.ErrorCtx(ctx, errutil.E(err).Debug("webrtc.NewTrackLocalStaticRTP"), "[voice_topic] NewTrackLocalStaticRTP")
			return
		}

		trackID := trackID(userID, remote)
		trackCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		publisher := r.getParticipant(userID)
		local := &localTrack{
			id:        trackID,
			owner:     userID,
			publisher: publisher,
			track:     track,
			cancel:    cancel,
			ssrc:      uint32(remote.SSRC()),
			kind:      remote.Kind(),
		}

		r.m.Lock()
		r.tracks[trackID] = local
		targets := make([]*participant, 0, len(r.participants))
		for _, participant := range r.participants {
			if participant.userID != userID {
				targets = append(targets, participant)
			}
		}
		r.m.Unlock()

		for _, target := range targets {
			if err := addTrackToParticipant(target, local); err != nil {
				logger.ErrorCtx(ctx, err, "[voice_topic] addTrackToParticipant")
				continue
			}
			r.renegotiate(ctx, target)
		}

		forwardRTP(trackCtx, remote, local)
		r.removeLocalTrack(trackID)
	}
}

func forwardRTP(ctx context.Context, remote *webrtc.TrackRemote, local *localTrack) {
	for {
		packet, _, err := remote.ReadRTP()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				logger.Error(errutil.E(err).Debug("remote.ReadRTP"), "[voice_topic] ReadRTP")
			}
			return
		}

		if err := local.track.WriteRTP(packet); err != nil {
			if !errors.Is(err, io.ErrClosedPipe) {
				logger.Error(errutil.E(err).Debug("local.track.WriteRTP"), "[voice_topic] WriteRTP")
			}
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (r *topic) removeLocalTrack(trackID string) {
	r.m.Lock()
	track := r.tracks[trackID]
	if track != nil {
		delete(r.tracks, trackID)
	}
	r.m.Unlock()

	if track == nil {
		return
	}
	track.cancel()
	r.removeTrackFromParticipants(trackID, track.owner, true)
}

func (r *topic) onICECandidate(userID domain.UserID) func(*webrtc.ICECandidate) {
	return func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		init := candidate.ToJSON()
		payload := event_desc.VoiceICECandidatePayload{
			Candidate:        init.Candidate,
			SDPMid:           init.SDPMid,
			SDPMLineIndex:    init.SDPMLineIndex,
			UsernameFragment: init.UsernameFragment,
		}
		r.enqueueICECandidate(userID, payload)
	}
}

func (r *topic) enqueueICECandidate(userID domain.UserID, candidate event_desc.VoiceICECandidatePayload) {
	participant := r.getParticipant(userID)
	if participant == nil {
		return
	}

	participant.m.Lock()
	participant.iceBatch = append(participant.iceBatch, candidate)
	if participant.iceBatchTimer == nil {
		participant.iceBatchTimer = time.AfterFunc(iceBatchDelay, func() {
			r.flushICECandidates(context.Background(), userID)
		})
	}
	participant.m.Unlock()
}

func (r *topic) flushICECandidates(ctx context.Context, userID domain.UserID) {
	participant := r.getParticipant(userID)
	if participant == nil {
		return
	}

	participant.m.Lock()
	candidates := participant.iceBatch
	participant.iceBatch = nil
	participant.iceBatchTimer = nil
	participant.m.Unlock()

	if len(candidates) == 0 {
		return
	}

	payload := event_desc.VoiceICECandidatesPayload{
		Candidates: candidates,
	}
	if err := r.sendEvent(ctx, userID, event_desc.TypeVoiceICECandidates, payload); err != nil {
		logger.ErrorCtx(ctx, err, "[voice_topic] send ICE candidates")
	}
}

func (r *topic) renegotiate(ctx context.Context, participant *participant) {
	participant.m.Lock()
	if participant.pc == nil ||
		!participant.remoteReady ||
		participant.renegotiating ||
		participant.pc.SignalingState() != webrtc.SignalingStateStable {
		participant.m.Unlock()
		return
	}
	participant.renegotiating = true
	participant.m.Unlock()

	offer, err := participant.pc.CreateOffer(nil)
	if err != nil {
		participant.m.Lock()
		participant.renegotiating = false
		participant.m.Unlock()
		logger.ErrorCtx(ctx, errutil.E(err).Debug("participant.pc.CreateOffer"), "[voice_topic] CreateOffer")
		return
	}
	if err := participant.pc.SetLocalDescription(offer); err != nil {
		participant.m.Lock()
		participant.renegotiating = false
		participant.m.Unlock()
		logger.ErrorCtx(ctx, errutil.E(err).Debug("participant.pc.SetLocalDescription"), "[voice_topic] SetLocalDescription")
		return
	}
	if err := r.sendSessionDescription(ctx, participant.userID, event_desc.TypeVoiceOffer, offer); err != nil {
		participant.m.Lock()
		participant.renegotiating = false
		participant.m.Unlock()
		logger.ErrorCtx(ctx, err, "[voice_topic] send renegotiation offer")
	}
}

func (r *topic) onConnectionStateChange(ctx context.Context, userID domain.UserID) func(webrtc.PeerConnectionState) {
	return func(state webrtc.PeerConnectionState) {
		switch state {
		case webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateClosed, webrtc.PeerConnectionStateDisconnected:
			if err := r.service.leaveLocal(ctx, r.id, userID); err != nil {
				logger.ErrorCtx(ctx, err, "[voice_topic] leaveLocal")
			}
		default:
		}
	}
}

func (r *topic) sendSessionDescription(ctx context.Context, userID domain.UserID, eventType event_desc.Type, desc webrtc.SessionDescription) error {
	payload := event_desc.VoiceSessionDescriptionPayload{
		Type: desc.Type.String(),
		SDP:  desc.SDP,
	}
	return r.sendEvent(ctx, userID, eventType, payload)
}

func (r *topic) sendEvent(ctx context.Context, userID domain.UserID, eventType event_desc.Type, payload any) error {
	event, err := render.Event(eventType, event_desc.Channel{
		Type: event_desc.ChannelVoiceTopic,
		ID:   r.id.I64(),
	}, payload)
	if err != nil {
		return errutil.E(err).Debug("render.Event")
	}
	if err := r.service.redis.PublishToUser(ctx, userID, event); err != nil {
		return errutil.E(err).Debug("r.service.redis.PublishToUser")
	}
	return nil
}

func trackID(userID domain.UserID, track *webrtc.TrackRemote) string {
	return strconv.FormatInt(userID.I64(), 10) + ":" + track.ID() + ":" + track.StreamID()
}

func readRTCP(sender *webrtc.RTPSender, track *localTrack) {
	for {
		packets, _, err := sender.ReadRTCP()
		if err != nil {
			return
		}
		for _, packet := range packets {
			switch packet.(type) {
			case *rtcp.PictureLossIndication, *rtcp.FullIntraRequest:
				requestKeyFrame(track)
			default:
			}
		}
	}
}

func (r *topic) requestParticipantKeyFrames(participant *participant) {
	r.m.RLock()
	tracks := make([]*localTrack, 0, len(r.tracks))
	for _, track := range r.tracks {
		if track.owner != participant.userID {
			tracks = append(tracks, track)
		}
	}
	r.m.RUnlock()

	participant.m.Lock()
	defer participant.m.Unlock()
	for _, track := range tracks {
		if _, ok := participant.senders[track.id]; ok {
			requestKeyFrameBurst(track)
		}
	}
}

func requestKeyFrameBurst(track *localTrack) {
	requestKeyFrame(track)
	time.AfterFunc(500*time.Millisecond, func() {
		requestKeyFrame(track)
	})
	time.AfterFunc(1500*time.Millisecond, func() {
		requestKeyFrame(track)
	})
}

func requestKeyFrame(track *localTrack) {
	if track.kind != webrtc.RTPCodecTypeVideo || track.publisher == nil || track.ssrc == 0 {
		return
	}

	track.publisher.m.Lock()
	pc := track.publisher.pc
	track.publisher.m.Unlock()
	if pc == nil {
		return
	}

	if err := pc.WriteRTCP([]rtcp.Packet{
		&rtcp.PictureLossIndication{MediaSSRC: track.ssrc},
	}); err != nil {
		logger.Error(errutil.E(err).Debug("publisher.pc.WriteRTCP"), "[voice_topic] request key frame")
	}
}
