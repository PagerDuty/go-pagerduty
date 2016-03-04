package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type IncidentManage struct {
}

func IncidentManageCommand() (cli.Command, error) {
	return &IncidentManage{}, nil
}

func (c *IncidentManage) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *IncidentManage) Synopsis() string {
	return "Acknowledge, resolve, escalate or reassign one or more incidents"
}

func (c *IncidentManage) Run(args []string) int {
	return 0
}
