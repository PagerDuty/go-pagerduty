package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type LogEntryShow struct {
}

func LogEntryShowCommand() (cli.Command, error) {
	return &LogEntryShow{}, nil
}

func (c *LogEntryShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *LogEntryShow) Synopsis() string {
	return "Get details for a specific incident log entry"
}

func (c *LogEntryShow) Run(args []string) int {
	return 0
}
