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

type MaintenanceWindowUpdate struct {
	Meta
}

func MaintenanceWindowUpdateCommand() (cli.Command, error) {
	return &MaintenanceWindowUpdate{}, nil
}

func (c *MaintenanceWindowUpdate) Help() string {
	helpText := `
	maintenance-window update <FILE> Update a maintenance window from json file
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowUpdate) Synopsis() string {
	return "Update an existing maintenance window"
}

func (c *MaintenanceWindowUpdate) Run(args []string) int {
	flags := c.Meta.FlagSet("maintenance-window update")
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
	var mw pagerduty.MaintenanceWindow
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

	if err := json.Unmarshal(f, &mw); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}

	log.Debugf("%#v", mw)
	if _, err := client.UpdateMaintenanceWindowWithContext(context.Background(), mw); err != nil {
		log.Error(err)
		return -1
	}

	return 0
}
