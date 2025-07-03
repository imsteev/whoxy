package clerk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ClerkIntegration struct {
}

func (c *ClerkIntegration) GetDestination(data []byte) (string, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	var event ClerkEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return "", err
	}

	ip := event.EventAttributes.HTTPRequest.ClientIP
	dest, err := redisClient.Get(context.Background(), ip).Result()
	if err != nil {
		return "", err
	}
	fmt.Println(dest)
	return dest, nil
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
