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

type ScheduleCreate struct {
	Meta
}

func ScheduleCreateCommand() (cli.Command, error) {
	return &ScheduleCreate{}, nil
}

func (c *ScheduleCreate) Help() string {
	helpText := `
	pd schedule create <FILE> Create a new schedule from json file
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *ScheduleCreate) Synopsis() string {
	return "Create a new on-call schedule"
}

func (c *ScheduleCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("schedule create")
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
	var s pagerduty.Schedule
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
	if err := decoder.Decode(&s); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", s)
	if _, err := client.CreateSchedule(s); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
