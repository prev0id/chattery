package voice

import "chattery/internal/domain"

type GetICEServersResponse struct {
	ICEServers []ICEServerResponse `json:"ice_servers"`
}

type ICEServerResponse struct {
	Username   string   `json:"username,omitzero"`
	Credential string   `json:"credential,omitzero"`
	URLs       []string `json:"urls"`
}

func convertGetICEServersResponse(iceServers []domain.VoiceICEServer) *GetICEServersResponse {
	response := &GetICEServersResponse{
		ICEServers: make([]ICEServerResponse, 0, len(iceServers)),
	}
	for _, iceServer := range iceServers {
		response.ICEServers = append(response.ICEServers, convertICEServerResponse(iceServer))
	}
	return response
}

func convertICEServerResponse(iceServer domain.VoiceICEServer) ICEServerResponse {
	return ICEServerResponse{
		URLs:       iceServer.URLs,
		Username:   iceServer.Username,
		Credential: iceServer.Credential,
	}
}
