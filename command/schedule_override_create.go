package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleOverrideCreate struct {
}

func ScheduleOverrideCreateCommand() (cli.Command, error) {
	return &ScheduleOverrideCreate{}, nil
}

func (c *ScheduleOverrideCreate) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleOverrideCreate) Synopsis() string {
	return "Create an override for a specific user"
}

func (c *ScheduleOverrideCreate) Run(args []string) int {
	return 0
}
