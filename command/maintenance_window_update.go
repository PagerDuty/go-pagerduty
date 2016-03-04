package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type MaintenanceWindowUpdate struct {
}

func MaintenanceWindowUpdateCommand() (cli.Command, error) {
	return &MaintenanceWindowUpdate{}, nil
}

func (c *MaintenanceWindowUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowUpdate) Synopsis() string {
	return "Update an existing maintenance window"
}

func (c *MaintenanceWindowUpdate) Run(args []string) int {
	return 0
}
