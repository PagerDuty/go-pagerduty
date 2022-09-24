package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleOverrideDelete struct {
}

func ScheduleOverrideDeleteCommand() (cli.Command, error) {
	return &ScheduleOverrideDelete{}, nil
}

func (c *ScheduleOverrideDelete) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleOverrideDelete) Synopsis() string {
	return "Remove an override"
}

func (c *ScheduleOverrideDelete) Run(args []string) int {
	return 0
}
