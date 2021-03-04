package main

import (
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
	if eps, err := client.ListEscalationPolicies(opts); err != nil {
		panic(err)
	} else {
		for _, p := range eps.EscalationPolicies {
			fmt.Println(p.Name)
		}
	}
}
