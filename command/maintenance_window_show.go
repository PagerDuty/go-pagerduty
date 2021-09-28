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

type MaintenanceWindowShow struct {
	Meta
}

func MaintenanceWindowShowCommand() (cli.Command, error) {
	return &MaintenanceWindowShow{}, nil
}

func (c *MaintenanceWindowShow) Help() string {
	helpText := `
	maintenance-window show Show a maintenance window

	Options:
		-id      The maintenance window ID
		-include Additional details to include
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowShow) Synopsis() string {
	return "Show an existing maintenance window"
}

func (c *MaintenanceWindowShow) Run(args []string) int {
	var includes []string
	var mwID string
	flags := c.Meta.FlagSet("maintenance-window show")
	flags.Usage = func() { fmt.Println(c.Help()) }
	flags.Var((*ArrayFlags)(&includes), "includes", "Additional details to include (can be specified multiple times)")
	flags.StringVar(&mwID, "id", "", "Maintenance window ID")

	if err := flags.Parse(args); err != nil {
		log.Error(err)
		return -1
	}

	if mwID == "" {
		log.Error("You must provide a maintenance window ID")
		return -1
	}

	if err := c.Meta.Setup(); err != nil {
		log.Error(err)
		return -1
	}

	client := c.Meta.Client()
	opts := pagerduty.GetMaintenanceWindowOptions{
		Includes: includes,
	}

	mw, err := client.GetMaintenanceWindowWithContext(context.Background(), mwID, opts)
	if err != nil {
		log.Error(err)
		return -1
	}

	data, err := yaml.Marshal(mw)
	if err != nil {
		log.Error(err)
		return -1
	}

	fmt.Println(string(data))
	return 0
}
