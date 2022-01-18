package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/PagerDuty/go-pagerduty/webhookv3"
)

const (
	secret = "lDQHScfUeXUKaQRNF+8XIiDKZ7XX3itBAYzwU0TARw8lJqRnkKl2iB1anSb0Z+IK"
)

func main() {
	http.HandleFunc("/webhook", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := webhookv3.VerifySignature(r, secret)
	if err != nil {
		switch err {
		case webhookv3.ErrNoValidSignatures:
			w.WriteHeader(http.StatusUnauthorized)

		case webhookv3.ErrMalformedBody, webhookv3.ErrMalformedHeader:
			w.WriteHeader(http.StatusBadRequest)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "%v", err)
		return
	}

	log.Infof("Received signed webhook")

	payload, err := webhookv3.ReadWebhookPayload(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("%v", err)
		return
	}

	event := payload.Event
	dataType, _ := event.GetEventDataValue("type")
	log.Infof("Event: %v, Event Type: %v, EventData Type: %v", event.ID, event.EventType, dataType)

	fmt.Fprint(w, "OK\n")
}
