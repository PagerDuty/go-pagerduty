package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ServiceIntegrationShow struct {
}

func ServiceIntegrationShowCommand() (cli.Command, error) {
	return &ServiceIntegrationShow{}, nil
}

func (c *ServiceIntegrationShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceIntegrationShow) Synopsis() string {
	return "Get details about an integration belonging to a service"
}

func (c *ServiceIntegrationShow) Run(args []string) int {
	return 0
}
