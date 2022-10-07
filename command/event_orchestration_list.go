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

type EventOrchestrationList struct {
	Meta
}

func (c *EventOrchestrationList) Help() string {
	helpText := `
	pd event-orchestration list List all of the existing event orchestrations

	Options:

		 -sort      Sort results by property (name:asc, name:desc, routes:asc, routes:desc, created_at:asc, created_at:desc)

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EventOrchestrationList) Synopsis() string {
	return "List all of the existing event orchestrations"
}

func EventOrchestrationListCommand() (cli.Command, error) {
	return &EventOrchestrationList{}, nil
}

func (c *EventOrchestrationList) Run(args []string) int {
	var sortBy string

	flags := c.Meta.FlagSet("event-orchestration list")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&sortBy, "sort", "", "Sort results by name, routes, or created_at (ascending or descending)")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	opts := pagerduty.ListOrchestrationsOptions{
		SortBy: sortBy,
	}
	if eps, err := client.ListOrchestrationsWithContext(context.Background(), opts); err != nil {
		log.Error(err)
		return -1
	} else {
		for i, p := range eps.Orchestrations {
			if i > 0 {
				fmt.Println("---")
			}
			data, err := yaml.Marshal(p)
			if err != nil {
				log.Error(err)
				return -1
			}
			fmt.Println(string(data))
		}
	}
	return 0
}
