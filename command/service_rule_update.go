package main

import (
	"strings"

	"github.com/mitchellh/cli"
)

type ServiceRuleUpdate struct {
}

func ServiceRuleUpdateCommand() (cli.Command, error) {
	return &ServiceRuleUpdate{}, nil
}

func (c *ServiceRuleUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceRuleUpdate) Synopsis() string {
	return "Get details about an integration belonging to a service"
}

func (c *ServiceRuleUpdate) Run(args []string) int {
	return 0
}
