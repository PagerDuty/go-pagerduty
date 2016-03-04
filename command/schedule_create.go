package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleCreate struct {
}

func ScheduleCreateCommand() (cli.Command, error) {
	return &ScheduleCreate{}, nil
}

func (c *ScheduleCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleCreate) Synopsis() string {
	return "Create a new on-call schedule"
}

func (c *ScheduleCreate) Run(args []string) int {
	return 0
}
