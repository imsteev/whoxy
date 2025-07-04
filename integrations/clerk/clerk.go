package clerk

import (
	"encoding/json"
	"net/http"
)

type ClerkIntegration struct {
}

func (c *ClerkIntegration) GetEventKey(r *http.Request) (string, error) {
	var event ClerkEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		return "", err
	}
	return event.EventAttributes.HTTPRequest.ClientIP, nil
}

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
