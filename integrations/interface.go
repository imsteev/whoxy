package integrations

import (
	"net/http"

	"github.com/imsteev/whoxy/integrations/clerk"
)

type Integration interface {
	// GetEventKey is used to derive a key from the request.
	// This key is typically used to identify the user that triggered a
	// webhook, which helps whoxy lookup where to forward the request.
	GetEventKey(*http.Request) (string, error)
}

var (
	_ Integration = (*clerk.ClerkIntegration)(nil)
	// Add more integrations here
	// Plugin system?
)
