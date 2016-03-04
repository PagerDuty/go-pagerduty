package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleShow struct {
}

func ScheduleShowCommand() (cli.Command, error) {
	return &ScheduleShow{}, nil
}

func (c *ScheduleShow) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleShow) Synopsis() string {
	return "Show detailed information about a schedule"
}

func (c *ScheduleShow) Run(args []string) int {
	return 0
}
