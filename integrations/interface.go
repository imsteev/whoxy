package integrations

import (
	"net/http"

	"github.com/imsteev/whoxy/integrations/clerk"
)

type Integration interface {
	// GetEventKey is used to get the event key from the request.
	// Event key is typically used to identify the user that triggered
	// a webhook, which then gets used to lookup where to forward the request.
	GetEventKey(*http.Request) (string, error)
}

var (
	_ Integration = (*clerk.ClerkIntegration)(nil)
)
