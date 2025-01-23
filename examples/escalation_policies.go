package main

import (
	"context"
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
)

var (
	subdomain string
	authtoken string
)

func ep() {
	var opts pagerduty.ListEscalationPoliciesOptions
	client := pagerduty.NewClient(authtoken)
	if eps, err := client.ListEscalationPoliciesWithContext(context.Background(), opts); err != nil {
		panic(err)
	} else {
		for _, p := range eps.EscalationPolicies {
			fmt.Println(p.Name)
		}
	}
}
