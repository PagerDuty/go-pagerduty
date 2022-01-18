package webhookv3

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/PagerDuty/go-pagerduty"
)

type OutboundEventData struct {
	Object map[string]interface{}
	// The raw json data for use in structured unmarshalling.
	RawData json.RawMessage
}

type OutboundEvent struct {
	ID           string                  `json:"id"`
	EventType    string                  `json:"event_type"`
	ResourceType string                  `json:"resource_type"`
	OccurredAt   string                  `json:"occurred_at"`
	Agent        *pagerduty.APIReference `json:"agent"`
	Data         *OutboundEventData      `json:"data"`
}

type WebhookPayload struct {
	Event OutboundEvent `json:"event"`
}

func (e *OutboundEventData) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.Object); err != nil {
		return err
	}

	if err := e.RawData.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}

func (oe *OutboundEvent) GetEventDataValue(keys ...string) (string, error) {
	return getDataValue(oe.Data.Object, keys)
}

func getDataValue(d map[string]interface{}, keys []string) (string, error) {
	key := keys[0]
	node := d[key]

	for k := 1; k < len(keys); k++ {
		key = keys[k]

		switch n := node.(type) {
		case []interface{}:
			intKey, err := strconv.Atoi(key)
			if err != nil {
				return "", fmt.Errorf("cannot identify array element with key '%s'", key)
			}
			node = n[intKey]
			continue
		case map[string]interface{}:
			node = n[key]
			continue
		default:
			break
		}
	}

	if node == nil {
		return "", fmt.Errorf("JSON does not have field '%s'", key)
	}

	return fmt.Sprintf("%v", node), nil
}
