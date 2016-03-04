package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type IncidentList struct {
}

func IncidentListCommand() (cli.Command, error) {
	return &IncidentList{}, nil
}

func (c *IncidentList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentList) Synopsis() string {
	return "List existing incidents"
}

func (c *IncidentList) Run(args []string) int {
	return 0
}
