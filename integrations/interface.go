package integrations

import "github.com/imsteev/whoxy/integrations/clerk"

type Integration interface {
	GetDestination([]byte) (string, error)
}

var (
	_ Integration = (*clerk.ClerkIntegration)(nil)
)
