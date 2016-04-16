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

type MaintenanceWindowCreate struct {
	Meta
}

func MaintenanceWindowCreateCommand() (cli.Command, error) {
	return &MaintenanceWindowCreate{}, nil
}

func (c *MaintenanceWindowCreate) Help() string {
	helpText := `
	pd maintenance-window create <FILE> Create a new maintenance window from json file
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowCreate) Synopsis() string {
	return "Create a new maintenance window for the specified services"
}

func (c *MaintenanceWindowCreate) Run(args []string) int {
	flags := c.Meta.FlagSet("maintenance-window create")
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
	var m pagerduty.MaintenanceWindow
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
	if err := decoder.Decode(&m); err != nil {
		log.Errorln("Failed to decode json. Error:", err)
		return -1
	}
	log.Debugf("%#v", m)
	if err := client.CreateMaintaienanceWindows(m); err != nil {
		log.Error(err)
		return -1
	}
	return 0
}
