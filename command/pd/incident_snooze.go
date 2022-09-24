package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type IncidentSnooze struct {
}

func IncidentSnoozeCommand() (cli.Command, error) {
	return &IncidentSnooze{}, nil
}

func (c *IncidentSnooze) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentSnooze) Synopsis() string {
	return "Snooze an incident"
}

func (c *IncidentSnooze) Run(args []string) int {
	return 0
}
