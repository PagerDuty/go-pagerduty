package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
)

var (
	subdomain string
	authtoken string
)

func main() {
	var opts pagerduty.ListEscalationPoliciesOptions
	client := pagerduty.NewClient(subdomain, authtoken)
	if eps, err := client.ListEscalationPolicies(opts); err != nil {
		panic(err)
	} else {
		for _, p := range eps.EscalationPolicies {
			fmt.Println(p.Name)
		}
	}
}
