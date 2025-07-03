package clerk

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type ClerkIntegration struct {
	Redis *redis.Client
}

func (c *ClerkIntegration) GetDestination(data []byte) (string, error) {
	var event ClerkEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return "", err
	}
	ip := event.EventAttributes.HTTPRequest.ClientIP
	return c.Redis.Get(context.Background(), ip).Result()
}
