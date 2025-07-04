package clerk

import (
	"context"
	"encoding/json"
	"net/http"
)

type ClerkIntegration struct {
}

func (c *ClerkIntegration) GetEventKey(r *http.Request) (string, error) {
	clonedReq := r.Clone(context.Background())
	var event ClerkEvent
	if err := json.NewDecoder(clonedReq.Body).Decode(&event); err != nil {
		return "", err
	}
	ip := event.EventAttributes.HTTPRequest.ClientIP
	return ip, nil
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
