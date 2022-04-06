package main

import (
	"github.com/mitchellh/cli"
	"strings"
)

type NotificationList struct {
}

func NotificationListCommand() (cli.Command, error) {
	return &NotificationList{}, nil
}

func (c *NotificationList) Help() string {
	helpText := `
	`
	return strings.TrimSpace(helpText)
}

func (c *NotificationList) Synopsis() string {
	return "List notifications for a given time range"
}

func (c *NotificationList) Run(args []string) int {
	return 0
}
