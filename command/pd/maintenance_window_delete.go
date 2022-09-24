package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

type MaintenanceWindowDelete struct {
	Meta
}

func MaintenanceWindowDeleteCommand() (cli.Command, error) {
	return &MaintenanceWindowDelete{}, nil
}

func (c *MaintenanceWindowDelete) Help() string {
	helpText := `
	maintenance-window delete Delete or end a maintenance window

	Options:
		-id      The maintenance window ID
	` + c.Meta.Help()
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowDelete) Synopsis() string {
	return "Delete an existing maintenance window if it's in the future, or end it if it's currently on-going"
}

func (c *MaintenanceWindowDelete) Run(args []string) int {
	var mwID string
	flags := c.Meta.FlagSet("maintenance-window delete")
	flags.Usage = func() { fmt.Println(c.Help()) }
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
	if err := client.DeleteMaintenanceWindowWithContext(context.Background(), mwID); err != nil {
		log.Error(err)
		return -1
	}

	return 0
}
