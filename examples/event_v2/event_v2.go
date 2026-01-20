package main

import (
	"context"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/google/uuid"
)

// TriggerAlert is a function that takes in a routing key and a message to
// create and send a PagerDuty alert. The service for where the PagerDuty
// alert is sent is determined at the RoutingKey level. The "message" is what
// the end user will see on the PagerDuty alert.
func TriggerAlert(routingKey string, message string) error {

	currentTime := time.Now()
	currentTimeAsString := currentTime.Format("2006-01-02T15:04:05.000Z")

	deduplicationKey := uuid.New()

	// Create a new V2 event to trigger an alert
	event := pagerduty.V2Event{
		RoutingKey: routingKey,
		Action:     "trigger",
		DedupKey:   deduplicationKey.String(),
		Payload: &pagerduty.V2Payload{
			Summary:   message,
			Source:    "Example Source",
			Severity:  "error", // critical, warning, error, info
			Timestamp: currentTimeAsString,
		},
	}

	// Send the V2 event to PagerDuty
	_, err := pagerduty.ManageEventWithContext(context.Background(), event)
	if err != nil {
		return err
	}

	return nil
}
