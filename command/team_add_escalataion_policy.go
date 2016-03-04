package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type TeamAddEscalationPolicy struct {
}

func TeamAddEscalationPolicyCommand() (cli.Command, error) {
	return &TeamAddEscalationPolicy{}, nil
}

func (c *TeamAddEscalationPolicy) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamAddEscalationPolicy) Synopsis() string {
	return "Add an escalation policy to a team"
}

func (c *TeamAddEscalationPolicy) Run(args []string) int {
	return 0
}
