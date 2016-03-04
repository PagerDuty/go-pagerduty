package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type EscalationPolicyUpdateUpdate struct {
}

func EscalationPolicyUpdateCommand() (cli.Command, error) {
	return &EscalationPolicyUpdateUpdate{}, nil
}

func (c *EscalationPolicyUpdateUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *EscalationPolicyUpdateUpdate) Synopsis() string {
	return "Updates an existing escalation policy and rules"
}

func (c *EscalationPolicyUpdateUpdate) Run(args []string) int {
	return 0
}
