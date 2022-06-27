package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type EventOrchestrationCreate struct {
	Meta
}

func EventOrchestrationCreateCommand() (cli.Command, error) {
	return &EventOrchestrationCreate{}, nil
}

func (c *EventOrchestrationCreate) Help() string {
	helpText := `
	pd event-orchestration create <FILE> Create a new event orchestration
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EventOrchestrationCreate) Synopsis() string {
	return "Creates a new event orchestration"
}

func (c *EventOrchestrationCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("event orchestration create")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	client := c.Meta.Client()
	var eo pagerduty.Orchestration
	if len(flags.Args()) != 1 {
		log.Error("Please specify input json file")
		return -1
	}
	log.Info("Input file is:", flags.Arg(0))
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&eo); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", eo)
	if _, err := client.CreateOrchestration(eo); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
