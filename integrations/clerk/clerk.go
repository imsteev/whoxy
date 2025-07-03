package clerk

import "encoding/json"

var (
	config = map[string]string{
		"2600:4040:a734:e600:a512:d9b:81d7:3615":  "https://thirty-poems-warn.loca.lt",
		"2600:4040:a734:e600:1c17:e974:645d:bb65": "https://thirty-poems-warn.loca.lt",
		"2600:4040:a734:e600:e8da:7217:5ccd:b4d2": "https://thirty-poems-warn.loca.lt",
		"fd07:b51a:cc66:0:a617:db5e:ab7:e9f1":     "https://thirty-poems-warn.loca.lt",
		"fd7a:115c:a1e0::5e01:8120":               "https://thirty-poems-warn.loca.lt",
		"2600:4040:a734:e600:9841:a8e3:167a:e21e": "https://thirty-poems-warn.loca.lt",
	}
)

type ClerkIntegration struct {
}

func (c *ClerkIntegration) GetDestination(data []byte) (string, error) {
	var event ClerkEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return "", err
	}
	ip := event.EventAttributes.HTTPRequest.ClientIP
	if dest, ok := config[ip]; ok {
		return dest, nil
	}
	return "", nil
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
