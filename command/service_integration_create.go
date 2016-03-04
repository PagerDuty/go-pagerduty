package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ServiceIntegrationCreate struct {
}

func ServiceIntegrationCreateCommand() (cli.Command, error) {
	return &ServiceIntegrationCreate{}, nil
}

func (c *ServiceIntegrationCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceIntegrationCreate) Synopsis() string {
	return "Create a new integration belonging to a service"
}

func (c *ServiceIntegrationCreate) Run(args []string) int {
	return 0
}
