package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type LogEntryList struct {
}

func LogEntryListCommand() (cli.Command, error) {
	return &LogEntryList{}, nil
}

func (c *LogEntryList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *LogEntryList) Synopsis() string {
	return "List all of the incident log entries across the entire account"
}

func (c *LogEntryList) Run(args []string) int {
	return 0
}
