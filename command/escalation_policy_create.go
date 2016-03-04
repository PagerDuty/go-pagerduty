package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type EscalationPolicyCreate struct {
}

func EscalationPolicyCreateCommand() (cli.Command, error) {
	return &EscalationPolicyCreate{}, nil
}

func (c *EscalationPolicyCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *EscalationPolicyCreate) Synopsis() string {
	return "Creates a new escalation policy"
}

func (c *EscalationPolicyCreate) Run(args []string) int {
	return 0
}
