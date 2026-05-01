package event_desc

import "encoding/json"

type Type string

const (
	TypePing               Type = "ping"
	TypePong               Type = "pong"
	TypeError              Type = "error"
	TypeJoin               Type = "join"
	TypeLeave              Type = "leave"
	TypeMessage            Type = "message"
	TypeVoiceOffer         Type = "voice_offer"
	TypeVoiceAnswer        Type = "voice_answer"
	TypeVoiceICECandidate  Type = "voice_ice_candidate"
	TypeVoiceICECandidates Type = "voice_ice_candidates"
	TypeVoiceState         Type = "voice_state"
	TypeVoiceJoined        Type = "voice_joined"
	TypeVoiceLeft          Type = "voice_left"
	TypeVoiceSignalAck     Type = "voice_signal_ack"
)

type Event struct {
	Type    Type            `json:"type"`
	Channel Channel         `json:"channel,omitzero"`
	Payload json.RawMessage `json:"payload,omitzero"`
}

type MessagePayload struct {
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Sender    UserInfo `json:"sender"`
	ID        int64    `json:"id"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

type VoiceParticipantPayload struct {
	NodeID  string   `json:"node_id"`
	Sender  UserInfo `json:"sender"`
	UserID  int64    `json:"user_id"`
	TopicID int64    `json:"topic_id"`
}

type VoiceStatePayload struct {
	Participants []VoiceParticipantPayload `json:"participants"`
	TopicID      int64                     `json:"topic_id"`
}

type VoiceSignalAckPayload struct {
	NodeID  string `json:"node_id"`
	TopicID int64  `json:"topic_id"`
}

type VoiceSessionDescriptionPayload struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

type VoiceICECandidatePayload struct {
	SDPMid           *string `json:"sdpMid,omitzero"`
	UsernameFragment *string `json:"usernameFragment,omitzero"`
	SDPMLineIndex    *uint16 `json:"sdpMLineIndex,omitzero"`
	Candidate        string  `json:"candidate"`
}

type VoiceICECandidatesPayload struct {
	Candidates []VoiceICECandidatePayload `json:"candidates"`
}

type UserInfo struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	ID       int64  `json:"id"`
}
