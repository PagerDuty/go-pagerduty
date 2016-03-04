package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type MaintenanceWindowShow struct {
}

func MaintenanceWindowShowCommand() (cli.Command, error) {
	return &MaintenanceWindowShow{}, nil
}

func (c *MaintenanceWindowShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowShow) Synopsis() string {
	return "Show an existing maintenance window"
}

func (c *MaintenanceWindowShow) Run(args []string) int {
	return 0
}
