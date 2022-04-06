package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type IncidentShow struct {
}

func IncidentShowCommand() (cli.Command, error) {
	return &IncidentShow{}, nil
}

func (c *IncidentShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentShow) Synopsis() string {
	return "Show detailed information about an incident"
}

func (c *IncidentShow) Run(args []string) int {
	return 0
}
