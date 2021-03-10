package main

import (
	"strings"

	"github.com/mitchellh/cli"
)

type ServiceRuleDelete struct {
}

func ServiceRuleDeleteCommand() (cli.Command, error) {
	return &ServiceRuleDelete{}, nil
}

func (c *ServiceRuleDelete) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ServiceRuleDelete) Synopsis() string {
	return "Get details about an integration belonging to a service"
}

func (c *ServiceRuleDelete) Run(args []string) int {
	return 0
}
