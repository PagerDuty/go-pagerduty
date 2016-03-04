package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type TeamUpdate struct {
}

func TeamUpdateCommand() (cli.Command, error) {
	return &TeamUpdate{}, nil
}

func (c *TeamUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *TeamUpdate) Synopsis() string {
	return "Update an existing team"
}

func (c *TeamUpdate) Run(args []string) int {
	return 0
}
