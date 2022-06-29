package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"

	"github.com/mitchellh/cli"
)

type EventOrchestrationUpdate struct {
	Meta
}

func EventOrchestrationUpdateCommand() (cli.Command, error) {
	return &EventOrchestrationUpdate{}, nil
}

func (c *EventOrchestrationUpdate) Help() string {
	helpText := `
	pd event-orchestration update <FILE> Update an event orchestration from json file
	Options:

		 -id

	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EventOrchestrationUpdate) Synopsis() string {
	return "Update an existing event orchestration"
}

func (c *EventOrchestrationUpdate) Run(args []string) int {
	var eoID string
	flags := c.Meta.FlagSet("event-orchestration update")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.StringVar(&eoID, "id", "", "Event orchestration id")
	if err := flags.Parse(args); err != nil {
		log.Error(err)
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
	var eo pagerduty.Orchestration
	if len(flags.Args()) != 1 {
		log.Error("Please specify input json file")
		return -1
	}

	log.Info("Input file is:", flags.Arg(0))
	f, err := ioutil.ReadFile(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}

	if err := json.Unmarshal(f, &eo); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}

	log.Debugf("%#v", eo)
	if _, err := client.UpdateOrchestrationWithContext(context.Background(), eoID, eo); err != nil {
		log.Error(err)
		return -1
	}

	return 0
}
