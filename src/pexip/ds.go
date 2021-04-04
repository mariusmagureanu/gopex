package pexip

// payload is a helper type for creating
// a request body payload.
type payload struct {
	DisplayName string `json:"display_name"`
}

// result is a helper type for parsing the
// response after a token request.
type result struct {
	Token              string `json:"token"`
	Expires            string `json:"expires"`
	ParticipantUUID    string `json:"participant_uuid"`
	DisplayName        string `json:"display_name"`
	AnalyticsEnabled   bool   `json:"analytics_enabled"`
	Role               string `json:"role"`
	ServiceType        string `json:"service_type"`
	ChatEnabled        bool   `json:"chat_enabled"`
	CurrentServiceType string `json:"current_service_type"`
}

// tokenResponse is a helper type for parsing the
// response after a token request.
type tokenResponse struct {
	Status string
	Result result
}
