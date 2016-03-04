package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type EscalationPolicyDelete struct {
}

func EscalationPolicyDeleteCommand() (cli.Command, error) {
	return &EscalationPolicyDelete{}, nil
}

func (c *EscalationPolicyDelete) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *EscalationPolicyDelete) Synopsis() string {
	return "Deletes an existing escalation policy and rules"
}

func (c *EscalationPolicyDelete) Run(args []string) int {
	return 0
}
