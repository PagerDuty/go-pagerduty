package main

import (
	"strings"

	"github.com/mitchellh/cli"
)

type ServiceRuleCreate struct {
}

func ServiceRuleCreateCommand() (cli.Command, error) {
	return &ServiceRuleCreate{}, nil
}

func (c *ServiceRuleCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceRuleCreate) Synopsis() string {
	return "Get details about an integration belonging to a service"
}

func (c *ServiceRuleCreate) Run(args []string) int {
	return 0
}
