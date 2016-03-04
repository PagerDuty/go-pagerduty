package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type MaintenanceWindowCreate struct {
}

func MaintenanceWindowCreateCommand() (cli.Command, error) {
	return &MaintenanceWindowCreate{}, nil
}

func (c *MaintenanceWindowCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *MaintenanceWindowCreate) Synopsis() string {
	return "Create a new maintenance window for the specified services"
}

func (c *MaintenanceWindowCreate) Run(args []string) int {
	return 0
}
