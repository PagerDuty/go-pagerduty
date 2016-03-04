package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type ScheduleList struct {
}

func ScheduleListCommand() (cli.Command, error) {
	return &ScheduleList{}, nil
}

func (c *ScheduleList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *ScheduleList) Synopsis() string {
	return "List the on-call schedules"
}

func (c *ScheduleList) Run(args []string) int {
	return 0
}
