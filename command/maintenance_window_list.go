package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type MaintenanceWindowList struct {
}

func MaintenanceWindowListCommand() (cli.Command, error) {
	return &MaintenanceWindowList{}, nil
}

func (c *MaintenanceWindowList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowList) Synopsis() string {
	return "List existing maintenance windows"
}

func (c *MaintenanceWindowList) Run(args []string) int {
	return 0
}
