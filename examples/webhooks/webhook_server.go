package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PagerDuty/go-pagerduty"
)

const (
	secret = "lDQHScfUeXUKaQRNF+8XIiDKZ7XX3itBAYzwU0TARw8lJqRnkKl2iB1anSb0Z+IK"
)

func main() {
	http.HandleFunc("/webhook", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	err := pagerduty.VerifySignature(body, r.Header.Get("X-PagerDuty-Signature"), secret)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "%v", err)
		return
	}

	fmt.Fprintf(w, "Received signed webhook")
}
