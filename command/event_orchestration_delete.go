package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type EventOrchestrationDelete struct {
	Meta
}

func EventOrchestrationDeleteCommand() (cli.Command, error) {
	return &EventOrchestrationDelete{}, nil
}

func (c *EventOrchestrationDelete) Help() string {
	helpText := `
	event-orchestration delete Delete an event orchestration

	Options:
		-id      The event orchestration ID
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EventOrchestrationDelete) Synopsis() string {
	return "Delete an existing event orchestration and its rules"
}

func (c *EventOrchestrationDelete) Run(args []string) int {
	var eoID string
	flags := c.Meta.FlagSet("event-orchestration delete")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&eoID, "id", "", "Event orchestration ID")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}

	if eoID == "" {
		log.Error("You must provide an event orchestration ID")
		return -1
	}

	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	if err := client.DeleteOrchestrationWithContext(context.Background(), eoID); err != nil {
		log.Error(err)
		return -1
	}

	return 0
}
