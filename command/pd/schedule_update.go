package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleUpdate struct {
}

func ScheduleUpdateCommand() (cli.Command, error) {
	return &ScheduleUpdate{}, nil
}

func (c *ScheduleUpdate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleUpdate) Synopsis() string {
	return "Update an existing on-call schedule"
}

func (c *ScheduleUpdate) Run(args []string) int {
	return 0
}
