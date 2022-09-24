package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleOverrideList struct {
}

func ScheduleOverrideListCommand() (cli.Command, error) {
	return &ScheduleOverrideList{}, nil
}

func (c *ScheduleOverrideList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleOverrideList) Synopsis() string {
	return "List overrides for a given time range"
}

func (c *ScheduleOverrideList) Run(args []string) int {
	return 0
}
