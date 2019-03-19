package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
)

type EventV2Manage struct {
	Meta
}

func EventV2ManageCommand() (cli.Command, error) {
	return &EventV2Manage{}, nil
}

func (c *EventV2Manage) Help() string {
	helpText := `
	pd event-v2 manage <FILE> Manage Events from json file
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *EventV2Manage) Synopsis() string {
	return "Create a New V2 Event"
}

func (c *EventV2Manage) Run(args []string) int {
	flags := c.Meta.FlagSet("event-v2 manage")
	flags.Usage = func() { fmt.Println(c.Help()) }
	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}
	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}
	var e pagerduty.V2Event
	if len(flags.Args()) != 1 {
		log.Error("Please specify input json file")
		return -1
	}
	log.Info("Input file is: ", flags.Arg(0))
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&e); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", e)
	if _, err := pagerduty.ManageEvent(e); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
