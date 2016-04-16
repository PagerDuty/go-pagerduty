package main

import (
	"encoding/json"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

type ScheduleOverrideCreate struct {
	Meta
}

func ScheduleOverrideCreateCommand() (cli.Command, error) {
	return &ScheduleOverrideCreate{}, nil
}

func (c *ScheduleOverrideCreate) Help() string {
	helpText := `
	pd schedule-override create <SERVICE-ID> <FILE> Create a schedule override from json file
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleOverrideCreate) Synopsis() string {
	return "Create an override for a specific user"
}

func (c *ScheduleOverrideCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("schedule-override create")
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
	var o pagerduty.Overrides
	if len(flags.Args()) != 2 {
		log.Error("Please specify input json file")
		return -1
	}
	log.Info("service id is:", flags.Arg(0))
	log.Info("Input file is:", flags.Arg(1))
	f, err := os.Open(flags.Arg(0))
	if err != nil {
		log.Error(err)
		return -1
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&o); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", o)
	if err := client.CreateOverride(flags.Arg(0), o); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
