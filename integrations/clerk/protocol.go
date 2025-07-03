package clerk

type ClerkEvent struct {
	Type            string          `json:"type"`
	EventAttributes EventAttributes `json:"event_attributes"`
}

type EventAttributes struct {
	HTTPRequest HTTPRequest `json:"http_request"`
}

type HTTPRequest struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}
