package webhookv3

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/PagerDuty/go-pagerduty"
)

// OutboundEventData is the unmarshalled data portion of the OutboundEvent.
type OutboundEventData struct {
	Object map[string]interface{}
	// The raw json data for use in structured unmarshalling.
	RawData json.RawMessage
}

// OutboundEvent represents the event that is delivered in a V3 Webhook Payload.
// See https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTkw-v3-overview#webhook-payload for more details.
type OutboundEvent struct {
	ID           string                  `json:"id"`
	EventType    string                  `json:"event_type"`
	ResourceType string                  `json:"resource_type"`
	OccurredAt   string                  `json:"occurred_at"`
	Agent        *pagerduty.APIReference `json:"agent"`
	Data         *OutboundEventData      `json:"data"`
}

// WebhookPayload represents the full object delivered as a result of V3 Webhook Subscriptions.
// See https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTkw-v3-overview#webhook-payload for more details.
type WebhookPayload struct {
	Event OutboundEvent `json:"event"`
}

// UnmarshalJSON is a custom unmarshaller used to produce OutboundEventData
// for further processing.
func (e *OutboundEventData) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.Object); err != nil {
		return err
	}

	if err := e.RawData.UnmarshalJSON(data); err != nil {
		return err
	}

	return nil
}

// GetEventDataValue returns a value from the e.Data object using the keys as a path
// or returns an error if the path does not point to a field.
//
// For example, `e.GetEventDataValue("type")` would return the `event.data.type` from a Webhook Payload.
// If the event type is `"incident"`, e.GetEventDataValue("priority", "id") would return the priority id.
// See the tests for additional examples.
func (e OutboundEvent) GetEventDataValue(keys ...string) (string, error) {
	return getDataValue(e.Data.Object, keys)
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
