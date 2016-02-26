package main

import (
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
)

/*
	Entrypoint for pd cli client based on the go bindings of API.
	good cli examples:
		https://github.com/Netflix-Skunkworks/go-jira
		https://github.com/docker/docker
		https://github.com/hashicorp/serf
*/

type Options struct {
	Authtoken string `short:"t" long:"token" description:"PagerDuty API Authentication token" required:"true"`
	Subdomain string `short:"d" long:"domain" description:"PagerDuty account name or DNS sub-domain" required:"true"`
	Loglevel  string `short:"l" long:"loglevel" description:"Log level" required:"false" default:"info"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err)
	}
	switch options.Loglevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatal("Unknown log level", options.Loglevel)
	}

	client := pagerduty.NewClient(options.Subdomain, options.Authtoken)
	var opts pagerduty.ListEscalationPoliciesOptions
	if eps, err := client.ListEscalationPolicies(opts); err != nil {
		log.Fatal(err)
	} else {
		for _, p := range eps.EscalationPolicies {
			log.Info(p)
		}
	}
}
