package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type EventOrchestrationShow struct {
	Meta
}

func (c *EventOrchestrationShow) Help() string {
	helpText := `
	pd event-orchestration show Show event orchestrations

	Options:

		 -id

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EventOrchestrationShow) Synopsis() string {
	return "Show information about an existing event orchestration and its rules"
}

func EventOrchestrationShowCommand() (cli.Command, error) {
	return &EventOrchestrationShow{}, nil
}

func (c *EventOrchestrationShow) Run(args []string) int {
	var eoID string
	flags := c.Meta.FlagSet("event-orchestration show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&eoID, "id", "", "Event orchestration id")
	if err := flags.Parse(args); err != nil {
		log.Errorln(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	if eoID == "" {
		log.Error("You must provide event orchestration id")
		return -1
	}
	client := c.Meta.Client()
	o := &pagerduty.GetOrchestrationOptions{}
	ep, err := client.GetOrchestrationWithContext(context.Background(), eoID, o)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err := yaml.Marshal(ep)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))
	fmt.Println("---")

	ro := &pagerduty.GetOrchestrationRouterOptions{}
	rules, err := client.GetOrchestrationRouterWithContext(context.Background(), eoID, ro)
	if err != nil {
		log.Error(err)
		return -1
	}
	data, err = yaml.Marshal(rules)
	if err != nil {
		log.Error(err)
		return -1
	}
	fmt.Println(string(data))

	return 0
}
